package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var s3Client *s3.S3
var region string

func init() {
	godotenv.Load()
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region = os.Getenv("AWS_REGION")

	// Configure AWS session to use the custom HTTP client
	awsConfig := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String("https://s3." + region + ".amazonaws.com"),
	}
	session, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Printf("error while creating session: %v\n", err)
		os.Exit(1)
	}

	// Create an S3 client with the configured session
	s3Client = s3.New(session)
}

func main() {
	// receive the name of the bucket as a flag parameter (default: bucketName)
	bucketNamePtr := flag.String("bucket", "bucketName", "the name of the bucket")
	flag.Parse()

	fmt.Printf("bucket name: %s\n", *bucketNamePtr)

	// list objects in the bucket
	listObjectsInput := &s3.ListObjectsInput{
		Bucket: aws.String(*bucketNamePtr),
	}

	// list objects in the bucket
	listObjectsOutput, err := s3Client.ListObjects(listObjectsInput)
	if err != nil {
		fmt.Printf("error while listing objects: %v\n", err)
		os.Exit(1)
	}

	// download each object in the bucket into a folder with the same name as the bucket
	bucketFolder := filepath.Join(".", *bucketNamePtr)
	for _, object := range listObjectsOutput.Contents {
		fmt.Printf("object: %s\n", aws.StringValue(object.Key))
		// check if object is a folder (ends with "/") if is a folder create it and continue to the next iteration
		if strings.HasSuffix(aws.StringValue(object.Key), "/") {
			fmt.Printf("creating folder: %s\n", aws.StringValue(object.Key))
			objectKey := strings.ReplaceAll(aws.StringValue(object.Key), "/", string(filepath.Separator))
			folderPath := filepath.Join(bucketFolder, objectKey)
			if err := os.MkdirAll(folderPath, 0755); err != nil {
				fmt.Printf("error while creating folder: %v\n", err)
				continue // Skip to the next iteration if there's an error
			}
		}


		fmt.Printf("downloading object: %s\n", aws.StringValue(object.Key))

		// download object
		getObjectInput := &s3.GetObjectInput{
			Bucket: aws.String(*bucketNamePtr),
			Key:    object.Key,
		}

		getObjectOutput, err := s3Client.GetObject(getObjectInput)
		if err != nil {
			fmt.Printf("error while getting object: %v\n", err)
			continue // Skip to the next iteration if there's an error
		}
		defer getObjectOutput.Body.Close() // Ensure closing the response body

		// create folder with the name of the bucket, if it doesn't exist
		if err := os.MkdirAll(bucketFolder, 0755); err != nil {
			fmt.Printf("error while creating folder: %v\n", err)
			continue // Skip to the next iteration if there's an error
		}

		// create file inside the bucket folder using sanitized object key
		objectKey := strings.ReplaceAll(aws.StringValue(object.Key), "/", string(filepath.Separator))
		filePath := filepath.Join(bucketFolder, objectKey)
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("error while creating file: %v\n", err)
			continue // Skip to the next iteration if there's an error
		}
		defer file.Close() // Ensure closing the file

		// write object to file
		_, err = io.Copy(file, getObjectOutput.Body)
		if err != nil {
			fmt.Printf("error while writing to file: %v\n", err)
			continue // Skip to the next iteration if there's an error
		}

		fmt.Printf("object downloaded successfully: %s\n", aws.StringValue(object.Key))
	}

	fmt.Println("done")
}