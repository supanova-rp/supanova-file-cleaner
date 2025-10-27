package main

import (
	"context"
	"fmt"
	"os"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supanova-rp/supanova-file-cleaner/config"
)

func main() {
	ctx := context.Background()

	err := run(ctx)
	if err != nil {
		fmt.Println("run failed:", err)
		os.Exit(1)
	}

	fmt.Println("app shutting down")
}

func run(ctx context.Context) error {
	cfg, err := config.ParseEnv()
	if err != nil {
		return fmt.Errorf("unable to parse env %v", err)
	}

	return listBucket(ctx, cfg)
}

func listBucket(ctx context.Context, cfg config.Config) error {
	awsConfig, err := aws_config.LoadDefaultConfig(
		ctx,
		aws_config.WithRegion(cfg.AWS.Region),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config %v", err)
	}

	client := s3.NewFromConfig(awsConfig)

	// TODO: paginate?
	resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.AWS.BucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to list items in bucket %v", err)
	}

	for _, item := range resp.Contents {
		fmt.Printf("Name: %s, Size: %d bytes\n", *item.Key, item.Size)
	}

	return nil
}
