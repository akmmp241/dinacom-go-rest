package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSClient struct {
	S3Client *s3.Client
	Uploader *manager.Uploader
}

func InitS3Client(cnf *Config) *AWSClient {
	region := cnf.Env.GetString("AWS_REGION")
	accessKey := cnf.Env.GetString("AWS_ACCESS_KEY_ID")
	secretKey := cnf.Env.GetString("AWS_SECRET_ACCESS_KEY")

	customCredentials := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		accessKey,
		secretKey,
		"",
	))

	s3Config, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(region),
		awsConfig.WithCredentialsProvider(customCredentials),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(s3Config)
	uploader := manager.NewUploader(s3Client)

	return &AWSClient{
		S3Client: s3Client,
		Uploader: uploader,
	}
}
