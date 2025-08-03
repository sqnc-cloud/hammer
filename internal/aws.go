package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var Uploader UploaderAPI

type UploaderAPI interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create session, %v", err))
	}

	Uploader = s3manager.NewUploader(sess)
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

func LoadAWSConfigFromEnv() (*AWSConfig, error) {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil, fmt.Errorf("environment variable AWS_ACCESS_KEY_ID not set")
	}

	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		return nil, fmt.Errorf("environment variable AWS_SECRET_ACCESS_KEY not set")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("environment variable AWS_REGION not set")
	}

	bucket := os.Getenv("AWS_BUCKET")
	if bucket == "" {
		return nil, fmt.Errorf("environment variable AWS_BUCKET not set")
	}

	return &AWSConfig{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
		Bucket:          bucket,
	}, nil
}

func UploadToS3(filePath, bucket, awsKey, awsSecret, region string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(awsKey, awsSecret, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create session, %v", err)
	}

	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}

	s3FileName := time.Now().Format(time.RFC3339) + ".zip"

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3FileName),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}
