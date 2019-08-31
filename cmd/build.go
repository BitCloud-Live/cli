package cmd

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/s3signer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yottab/cli/config"
	ybApi "github.com/yottab/proto-api/proto"
)

const (
	s3UriFormat            = "http://%s/%s/%s"           // http://endpoint/bucketName/objectName
	s3DefaulteRegion       = "us-east-1"                 //
	s3DefaultesessionToken = ""                          //
	imageLogIDFormat       = "%s:%s"                     // imageName:imageTag
	archiveNameFormat      = "repository_%s_archive.zip" //
)

func imageBuild(cmd *cobra.Command, args []string) {
	baseFolder := cmd.Flag("path").Value.String() // Repository path
	imageTag := cmd.Flag("tag").Value.String()    // Repository tag
	imageName := cmd.Flag("name").Value.String()  // Repository name
	zipName := fmt.Sprintf(archiveNameFormat, imageName)

	checkExistDockerfile(baseFolder)
	zipPath, err := zipFolder(baseFolder, zipName)
	uiCheckErr("Could not Archive the Folder", err)
	log.Printf("archive folder at [%s]", zipPath)

	zipArchiveURL, err := s3SendArchive(zipPath, zipName)
	uiCheckErr("Could not Save the Archive at s3.YOTTAb.io", err)
	log.Printf("Upload Archive at [%s]", zipArchiveURL)

	client := grpcConnect()
	defer client.Close()
	req := new(ybApi.ImgBuildReq)
	req.RepositoryName = imageName
	req.RepositoryTag = imageTag
	req.ZipArchiveURL = zipArchiveURL
	_, err = client.V2().ImgBuild(client.Context(), req)
	uiCheckErr("Could not Build the Repository", err)
	getBuildLog(imageName, imageTag)
}

func getBuildLog(imageName, imageTag string) {
	id := getRequestIdentity(
		fmt.Sprintf(imageLogIDFormat, imageName, imageTag))
	client := grpcConnect()
	defer client.Close()
	logClient, err := client.V2().ImgBuildLog(context.Background(), id)
	uiCheckErr("Could not Get Application log", err)
	uiImageLog(logClient)
}

func initializeS3ArchiveBucket(minioClient *minio.Client, bucketName string) (err error) {
	// Check to see if we already own this bucket
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Printf("Err: check Bucket Exists; Bucket:%s, Err:%v", bucketName, err)
		return
	} else if !exists {
		// Make a new bucket.
		if err = minioClient.MakeBucket(bucketName, s3DefaulteRegion); err != nil {
			log.Printf("Err: Make Bucket %s, Err:%v", bucketName, err)
			return
		}
	}
	return nil
}

func zipBufferIO(zipFilePath, objectName string) (bodyBuf *bytes.Buffer, err error) {
	file, err := os.Open(zipFilePath)
	if err != nil {
		log.Printf("Err: Accesse to Archive at [%s], Err:%v", zipFilePath, err)
		return
	}
	defer file.Close()

	bodyBuf = new(bytes.Buffer)
	_, err = bodyBuf.ReadFrom(file)
	if err != nil {
		log.Printf("Err: Buffer Reader From file [%s], Err:%v", zipFilePath, err)
	}
	return
}

func s3PutObject(minioClient *minio.Client, zipFilePath, accessKeyID, secretAccessKey, endpoint, bucketName, objectName string) (err error) {
	/*/
		TODO: raplace when bakend is minio
		TODO: remove extra func.arg

		n, err := minioClient.FPutObject(bucketName, objectName, zipFilePath, minio.PutObjectOptions{
			ContentType: "x-www-form-urlencoded",
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully uploaded bytes: ", n)
	/*/

	// conver Zip to IO.Writer
	bodyBuf, err := zipBufferIO(zipFilePath, objectName)
	if err != nil {
		log.Printf("Err: zipBufferIO() Path:[%s], Err:%v", zipFilePath, err)
		return
	}

	// PUT zip to s3
	s3ObjectNameURI := fmt.Sprintf(s3UriFormat, endpoint, bucketName, objectName)
	client := &http.Client{}

	s3Req, err := http.NewRequest(http.MethodPut, s3ObjectNameURI, bodyBuf)
	if err != nil {
		log.Printf("Err: http PUT Request URI:[%s] Err:%v", s3ObjectNameURI, err)
		return
	}
	s3Req = s3signer.SignV4(*s3Req, accessKeyID, secretAccessKey, s3DefaultesessionToken, s3DefaulteRegion)
	s3Req.Header.Set("Content-Type", "application/octet-stream")
	_, err = client.Do(s3Req)
	if err != nil {
		log.Printf("Err: Send file URI:[%s] Err:%v", s3ObjectNameURI, err)
	}
	return
}

func s3SendArchive(zipFilePath, objectName string) (uri string, err error) {
	var (
		endpoint        = "s3.YOTTAb.io"
		bucketName      = fmt.Sprintf("yb--build-archive--%s", viper.GetString(config.KEY_USER))
		accessKeyID     = viper.GetString(config.KEY_TOKEN)
		secretAccessKey = " "
		useSSL          = false
	)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Printf("Err: Initialize s3 client, Err:%v", err)
		return
	}

	// Initialize Archive bucket.
	if err = initializeS3ArchiveBucket(minioClient, bucketName); err != nil {
		log.Printf("Err: Initialize S3 Archive bucket, Err:%v", err)
		return
	}

	// save Archive File
	if err = s3PutObject(minioClient, zipFilePath, accessKeyID, secretAccessKey, endpoint, bucketName, objectName); err != nil {
		log.Printf("Err: Put S3 Archive bucket, Err:%v", err)
		return
	}

	// Genrate URL for Download
	uriObj, err := minioClient.PresignedGetObject(bucketName, objectName, time.Hour*16, nil)
	if err != nil {
		log.Printf("Err: %s Put Object to Bucket:%s, Path:%s, Name:%s, Err:%v", endpoint, bucketName, zipFilePath, objectName, err)
		return
	}
	return uriObj.String(), nil
}

func zipFolder(baseFolder, archFileName string) (archivePath string, err error) {
	archivePath = fmt.Sprintf("%s%s", baseFolder, archFileName)

	// Get a Buffer to Write To
	outFile, err := os.Create(archivePath)
	uiCheckErr("Could not create file", err)
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	zipAddFiles(w, baseFolder, "", archFileName)

	if err != nil {
		os.Remove(archivePath)
		log.Fatalf("%v", err)
	}

	// Make sure to check the error on Close.
	if err = w.Close(); err != nil {
		os.Remove(archivePath)
		log.Fatalf("%v", err)
	}

	return
}

func checkExistDockerfile(basePath string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatalf("ioutil.ReadDir Err: %v", err)
	}

	for _, file := range files {
		if file.Name() == "Dockerfile" {
			return
		}
	}
	log.Fatalf("Err: can find Dockerfile at [%s]", basePath)
}

func zipAddFiles(w *zip.Writer, basePath, baseInZip, archFileName string) (err error) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Printf("ioutil.ReadDir Err: %v", err)
		return err
	}

	for _, file := range files {
		// archive file created at rootPath, dont add it
		if len(baseInZip) == 0 && file.Name() == archFileName {
			continue
		}

		if file.IsDir() {
			// Recurse
			newBase := fmt.Sprintf("%s%s%c", basePath, file.Name(), os.PathSeparator)
			newBaseInZip := fmt.Sprintf("%s%s%c", baseInZip, file.Name(), os.PathSeparator)

			if err = zipAddFiles(w, newBase, newBaseInZip, archFileName); err != nil {
				log.Printf("zipAddFiles basePath: %s, file:%s, Err: %v", basePath, file.Name(), err)
				return err
			}
		} else {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				log.Printf("zip ioutil.ReadFile Err: %v", err)
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				log.Printf("zip.Writer.Create Err: %v", err)
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				log.Printf("zip.Writer.Write Err: %v", err)
				return err
			}
		}
	}
	return
}
