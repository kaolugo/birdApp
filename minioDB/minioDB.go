package minioDB

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/joho/godotenv"
)

// make a minio struct ?? I think ??
type MinioDB struct {
	Client *s3.S3
}

// functions

//create a new Minio database !!
func NewImageDatabase(address string) (*MinioDB, *string, error) {
	bucket := aws.String("imageBucket")

	// let's use environment variabes
	err1 := godotenv.Load(".env")

	if err1 != nil {
		log.Fatalf("Error loading .env file")
	}

	access := os.Getenv("ACCESS")
	secret := os.Getenv("SECRET")

	// configure to use MinIO server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(access, secret, ""),
		Endpoint:         aws.String(address),
		Region:           aws.String("ap-northeast-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	cparams := &s3.CreateBucketInput{
		Bucket: bucket,
	}

	// create a new bucket using the CreateBucket call
	_, err := s3Client.CreateBucket(cparams)
	if err != nil {
		// message from an error
		fmt.Println(err.Error())
		return nil, bucket, err
	}

	// return the s3 client as well as the bucket name string
	return &MinioDB{
		Client: s3Client,
	}, bucket, nil
}
