package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// func main() {
// 	// Load the Shared AWS Configuration (~/.aws/config)
// 	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKIAY2JPXHV4BVEVAGA3", "uzF5YBGTZ02PysfxjwN4hRwutBhKLoHGAz+e+VPs", "")))
// 	if err != nil {
// 		log.Fatal("Error 1 - ", err)
// 	}

// 	// Create an Amazon S3 service client
// 	client := s3.NewFromConfig(cfg)

// 	// Get the first page of results for ListObjectsV2 for a bucket
// 	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
// 		Bucket: aws.String("sandeep-buckets"),
// 	})
// 	if err != nil {
// 		log.Fatal("Error 2 - ", err)
// 	}

// 	log.Println("first page results:")
// 	for _, object := range output.Contents {
// 		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
// 	}
// }

var (
	s3session *s3.S3
)

const (
	BUCKET_NAME = "sandeep-buckets-3"
	REGIN       = "ap-south-1"
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGIN),
	})))
}

func ListBuckets() (resp *s3.ListBucketsOutput) {
	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	return resp
}

func CreateBucket(bucketName string) (resp *s3.CreateBucketOutput) {
	resp, err := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGIN),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket name is already in use!")
				panic(err)

			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket exists and is owned by you.")
				return nil

			default:
				panic(err)
			}
		}
	}
	return resp
}

func UploadObject(filename string) (resp *s3.PutObjectOutput) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("Uploading:", filename)
	resp, err = s3session.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(strings.Split(filename, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func ListObjects() (resp *s3.ListObjectsV2Output) {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BUCKET_NAME),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func GetObject(filename string) {
	fmt.Println("Downloading: ", filename)

	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}
}

func DeleteObject(filename string) (resp *s3.DeleteObjectOutput) {
	fmt.Println("Deleting: ", filename)
	resp, err := s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func main() {
	fmt.Println(ListBuckets())
	// fmt.Println(CreateBucket(BUCKET_NAME))

	folder := "files"
	files, _ := ioutil.ReadDir(folder)
	fmt.Println(files)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			UploadObject(folder + "/" + file.Name())
		}
	}

	fmt.Println(ListObjects())

	for _, object := range ListObjects().Contents {
		GetObject(*object.Key)
		DeleteObject(*object.Key)
	}

	fmt.Println(ListObjects())
}
