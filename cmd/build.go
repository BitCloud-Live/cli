package cmd

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yottab/cli/config"
	ybApi "github.com/yottab/proto-api/proto"
)

func repositoryBuild(cmd *cobra.Command, args []string) {
	req := new(ybApi.RepositoryBuildReq)
	repositoryPath := cmd.Flag("path").Value.String()  // Repository path
	req.RepositoryTag = cmd.Flag("tag").Value.String() // Repository tag
	repositoryName := cmd.Flag("name").Value.String()  // Repository name
	req.RepositoryName = repositoryName
	archFileName := fmt.Sprintf("repository_%s.zip", repositoryName)

	zipPath, err := zipArchive(repositoryPath, archFileName)
	uiCheckErr("Could not Archive the Repository", err)

	req.ZipArchiveURL, err = s3SendArchive(zipPath, archFileName)
	uiCheckErr("Could not Save the Archive at s3.YOTTAb.io", err)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().RepositoryBuild(client.Context(), req)

	uiCheckErr("Could not Build the Repository", err)
	log.Println(res.Log)
}

func s3SendArchive(filePath, objectName string) (uri string, err error) {
	var (
		endpoint        = "s3.YOTTAb.io"
		bucketName      = fmt.Sprintf("Yb--RepositoryBuildArchive--%s", viper.GetString(config.KEY_USER))
		accessKeyID     = viper.GetString(config.KEY_TOKEN)
		secretAccessKey = " "
		location        = "us-east-1"
		useSSL          = false
		contentType     = "application/zip"
	)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Printf("Initialize s3 client, Err:%v", err)
		return
	}

	// Check to see if we already own this bucket
	exists, errBucketExists := minioClient.BucketExists(bucketName)
	if errBucketExists == nil && !exists {
		// Make a new bucket.
		if err = minioClient.MakeBucket(bucketName, location); err != nil {
			log.Printf("%s Make Bucket %s, Err:%v", endpoint, bucketName, err)
			return
		}
	}

	// Upload the zip file
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("%s Put Object to Bucket %s, Err:%v", endpoint, bucketName, err)
		return
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName), nil
}

func zipArchive(baseFolder, archFileName string) (archivePath string, err error) {
	archivePath = fmt.Sprintf("%s%v%s", archivePath, os.PathSeparator, archFileName)

	// Get a Buffer to Write To
	outFile, err := os.Create(archivePath)
	uiCheckErr("Could not create file", err)
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

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

func addFiles(w *zip.Writer, basePath, baseInZip string) (err error) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Printf("ioutil.ReadDir Err: %v", err)
		return err
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				log.Printf("ioutil.ReadFile Err: %v", err)
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
		} else if file.IsDir() {
			// Recurse
			newBase := fmt.Sprintf("%s%s%v", basePath, file.Name(), os.PathSeparator)
			newBaseInZip := fmt.Sprintf("%s%v", file.Name(), os.PathSeparator)
			fmt.Println("Recursing and Adding SubDir: " + newBase)
			err = addFiles(w, newBase, newBaseInZip)
			if err != nil {
				break
			}
		}
	}
	return
}
