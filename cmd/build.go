package cmd

import (
	"archive/zip"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	ignore "github.com/codeskyblue/dockerignore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yottab/cli/config"
	ybApi "github.com/yottab/proto-api/proto"
)

const (
	imageLogIDFormat  = "%s:%s"                     // imageName:imageTag
	archiveNameFormat = "repository_%s_archive.zip" //
	bucketNameFormat  = "yottab-bucket-build-%s"    //
)

var (
	ignoreFile = ".dockerignore" // TODO get by EnvVar
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
	log.Print("Successful Uploading Archive.")

	client := grpcConnect()
	defer client.Close()
	req := new(ybApi.ImgBuildReq)
	req.RepositoryName = imageName
	req.RepositoryTag = imageTag
	req.ZipArchiveURL = zipArchiveURL
	_, err = client.V2().ImgBuild(client.Context(), req)
	uiCheckErr("Could not Build the Repository", err)
	log.Print("Build started!")
	log.Print("Waiting for builder log to get ready...")
	time.Sleep(20 * time.Second)
	getBuildLog(imageName, imageTag)
}

func imageBuildLog(cmd *cobra.Command, args []string) {
	imageTag := cmd.Flag("tag").Value.String()   // Repository tag
	imageName := cmd.Flag("name").Value.String() // Repository name
	getBuildLog(imageName, imageTag)
}

func getBuildLog(imageName, imageTag string) {
	id := getRequestIdentity(
		fmt.Sprintf(imageLogIDFormat, imageName, imageTag))
	client := grpcConnect()
	defer client.Close()
	logClient, err := client.V2().ImgBuildLog(context.Background(), id)
	uiCheckErr(fmt.Sprintf("Could not get build log right now!\nTry again in a few soconds using $yb build-log --name=%s --tag=%s", imageName, imageTag), err)
	uiImageLog(logClient)

}

func s3SendArchive(zipFilePath, objectName string) (uri string, err error) {
	var bucketName = fmt.Sprintf(bucketNameFormat, viper.GetString(config.KEY_USER))

	// Initialize minio client object.
	minioClient := initializeObjectStore()

	// Initialize Archive bucket.
	if err = initializeS3ArchiveBucket(minioClient, bucketName); err != nil {
		log.Printf("Err: Initialize S3 Archive bucket, Err:%v", err)
		return
	}

	// save Archive File
	if err = s3PutObject(minioClient, zipFilePath, bucketName, objectName); err != nil {
		log.Printf("Err: Put S3 Archive bucket, Err:%v", err)
		return
	}

	// Genrate URL for Download
	uriObj, err := minioClient.PresignedGetObject(bucketName, objectName, time.Hour*16, nil)
	if err != nil {
		log.Printf("Err: Put Object to Bucket:%s, Path:%s, Name:%s, Err:%v", bucketName, zipFilePath, objectName, err)
		return
	}
	return uriObj.String(), nil
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

func zipFolder(baseFolder, archFileName string) (archivePath string, err error) {
	ignorePatterns := readIgnorePatterns()
	archivePath = fmt.Sprintf("%s%s", baseFolder, archFileName)

	// Get a Buffer to Write To
	outFile, err := os.Create(archivePath)
	uiCheckErr("Could not create file", err)
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	zipAddFiles(w, baseFolder, "", archFileName, ignorePatterns)

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

func zipAddFiles(w *zip.Writer, basePath, baseInZip, archFileName string, ignorePatterns []string) (err error) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Printf("ioutil.ReadDir Err: %v", err)
		return err
	}

	for _, file := range files {
		// if path in ignorePatterns, Skip path
		if checkIgnoreMatches(
			baseInZip+file.Name(), // File path in project
			ignorePatterns) {
			continue
		}

		// archive file created at rootPath, dont add it
		if len(baseInZip) == 0 && (file.Name() == archFileName || file.Name() == ignoreFile) {
			continue
		}

		if file.IsDir() {
			// Recurse
			newBase := fmt.Sprintf("%s%s%c", basePath, file.Name(), os.PathSeparator)
			newBaseInZip := fmt.Sprintf("%s%s%c", baseInZip, file.Name(), os.PathSeparator)

			if err = zipAddFiles(w, newBase, newBaseInZip, archFileName, ignorePatterns); err != nil {
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

func readIgnorePatterns() (patterns []string) {
	patterns, err := ignore.ReadIgnoreFile(ignoreFile)
	if os.IsNotExist(err) {
		return []string{}
	}
	uiCheckErr("Read '.dockerignore' file", err)
	log.Printf("Successfully reading data from of file [%s]", ignoreFile)
	return
}

func checkIgnoreMatches(path string, patterns []string) (isSkip bool) {
	isSkip, err := ignore.Matches(path, patterns)
	uiCheckErr("DockerIgnore check for path ["+path+"]", err)
	return
}
