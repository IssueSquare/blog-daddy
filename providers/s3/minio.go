package s3

import (
	"github.com/minio/minio-go"
	"log"
)

type (
	S3Repo interface {
		Upload(bucketName, filePath string) error
		CreateBucket(bucketName string) error
	}

	S3 struct {
		endpoint  string
		accessKey string
		secretKey string
	}
)

func NewS3Provider(endpoint, accessKey, secretKey string) *S3 {
	s3 := new(S3)
	s3.endpoint = endpoint
	s3.accessKey = accessKey
	s3.secretKey = secretKey
	return s3
}

func (s *S3) Upload(bucketName, filePath string) error {
	// Initialize minio client object.
	minioClient, err := minio.New(s.endpoint, s.accessKey, s.secretKey, true)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	// Upload the zip file
	objectName := "index.html"
	contentType := "text/html"

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	return nil
}

func (s *S3) CreateBucket(bucketName string) error {

	// Initialize minio client object.
	minioClient, err := minio.New(s.endpoint, s.accessKey, s.secretKey, false)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("%#v\n", minioClient) // minioClient is now setup

	//Default
	location := "us-east-1"
	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	return nil
}
