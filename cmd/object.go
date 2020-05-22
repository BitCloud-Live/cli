package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func checkIsNotDir(srcPath string) {
	file, err := os.Open(srcPath)
	uiCheckErr("Could not Accesse to File", err)
	defer file.Close()

	fs, err := file.Stat()
	uiCheckErr("Could not Get the Stat if File", err)
	if fs.IsDir() {
		log.Fatal("Could not Copy Folder")
	}
}

func detectBucket(path string) (bucket, object string) {
	if strings.Index(strings.ToLower(path), "yottab.io") > -1 {
		log.Fatal("Error: bad format of Object path, enter [ Bucket_Name/Object_Path ]")
	}

	splitPath := strings.Split(path, "/")
	if len(splitPath) < 2 {
		log.Fatal("Error: bad format of Object path, enter [ Bucket_Name/Object_Path ]")
	}

	bucket = splitPath[0]
	if len(bucket) < 1 {
		log.Fatal("Error: bad format of Object path, enter [ Bucket_Name/Object_Path ]")
	}

	object = path[1+len(bucket):]
	return
}

//TODO
func objectLs(cmd *cobra.Command, args []string) {}

func objectCp(cmd *cobra.Command, args []string) {
	var (
		srcPath                = getCliRequiredArg(args, 0)
		DesPath                = getCliRequiredArg(args, 1)
		bucketName, objectName = detectBucket(DesPath)
	)

	// Only support Copy Single File
	checkIsNotDir(srcPath)

	// Initialize minio client object.
	minioClient := initializeObjectStore()

	// Initialize Archive bucket.
	if err := initializeS3ArchiveBucket(minioClient, bucketName); err != nil {
		log.Printf("Err: Initialize S3 Archive bucket, Err:%v", err)
		return
	}

	// save Archive File
	if err := s3PutObject(minioClient, srcPath, bucketName, objectName); err != nil {
		log.Printf("Err: Put S3 Archive bucket, Err:%v", err)
		return
	}
}

func objectRm(cmd *cobra.Command, args []string) {
	var (
		objectPath             = getCliRequiredArg(args, 0)
		bucketName, objectName = detectBucket(objectPath)
	)

	// Initialize minio client object.
	minioClient := initializeObjectStore()

	err := minioClient.RemoveObject(bucketName, objectName)
	uiCheckErr("Could not Remove Object", err)

	log.Print("Successful delete")
}

func bucketCreate(cmd *cobra.Command, args []string) {
	bucketName := getCliRequiredArg(args, 0)

	// Initialize minio client object.
	minioClient := initializeObjectStore()

	err := minioClient.MakeBucket(bucketName, s3DefaulteRegion)
	uiCheckErr("Could not Create Object", err)
	fmt.Println("Successfully created Bucket.")
}

func bucketList(cmd *cobra.Command, args []string) {
	// Initialize minio client object.
	minioClient := initializeObjectStore()

	buckets, err := minioClient.ListBuckets()
	uiCheckErr("Could not Buckets List", err)

	for i, bucket := range buckets {
		log.Printf("%3d. [%s]  \tName: %s",
			i+1,
			bucket.CreationDate.Format("Mon Jan _2 15:04:05 2006"),
			bucket.Name)
	}
}
