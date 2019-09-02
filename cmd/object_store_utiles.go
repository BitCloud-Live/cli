package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/s3signer"
	"github.com/spf13/viper"
	"github.com/yottab/cli/config"
)

const (
	s3UriFormat            = "http://%s/%s/%s" // http://endpoint/bucketName/objectName
	s3DefaulteRegion       = "us-east-1"       //
	s3DefaultesessionToken = ""                //
)

var (
	s3Endpoint        = "s3.YOTTAb.io"                    // TODO get by EVar
	s3AccessKeyID     = viper.GetString(config.KEY_TOKEN) // TODO get by EVar
	s3SecretAccessKey = " "                               // TODO get by EVar
	s3UseSSL          = false                             // TODO get by EVar
)

// Initialize minio client object.
func initializeObjectStore() (minioClient *minio.Client, err error) {
	minioClient, err = minio.New(s3Endpoint, s3AccessKeyID, s3SecretAccessKey, s3UseSSL)
	if err != nil {
		log.Printf("Err: Initialize s3 client, Err:%v", err)
	}

	return
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

// TODO kill me (when backend is minio)
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

// TODO edit code when backend is minio
func s3PutObject(minioClient *minio.Client, zipFilePath, bucketName, objectName string) (err error) {
	/*/
		TODO: raplace when backend is minio

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
	s3ObjectNameURI := fmt.Sprintf(s3UriFormat, s3Endpoint, bucketName, objectName)
	client := &http.Client{}

	s3Req, err := http.NewRequest(http.MethodPut, s3ObjectNameURI, bodyBuf)
	if err != nil {
		log.Printf("Err: http PUT Request URI:[%s] Err:%v", s3ObjectNameURI, err)
		return
	}
	s3Req = s3signer.SignV4(*s3Req, s3AccessKeyID, s3SecretAccessKey, s3DefaultesessionToken, s3DefaulteRegion)
	s3Req.Header.Set("Content-Type", "application/octet-stream")
	_, err = client.Do(s3Req)
	if err != nil {
		log.Printf("Err: Send file URI:[%s] Err:%v", s3ObjectNameURI, err)
	}
	return
}
