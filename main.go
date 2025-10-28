package main

import (
	"context"
	"fmt"
	"os"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supanova-rp/supanova-file-cleaner/internal/config"
	"github.com/supanova-rp/supanova-file-cleaner/internal/store"
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
		return fmt.Errorf("unable to parse env: %v", err)
	}

	db, err := store.NewStore(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer db.Close()

	err = listBucket(ctx, cfg.AWS)
	if err != nil {
		return fmt.Errorf("failed to list bucket: %v", err)
	}

	videos, err := db.Queries.GetVideos(ctx)
	if err != nil {
		return fmt.Errorf("failed to get videos: %v", err)
	}

	fmt.Printf("\n%+v\n", videos)

	// TODO: Best way to make the app block?
	select {}
}

func listBucket(ctx context.Context, cfg config.AWSConfig) error {
	awsConfig, err := aws_config.LoadDefaultConfig(
		ctx,
		aws_config.WithRegion(cfg.Region),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config %v", err)
	}

	client := s3.NewFromConfig(awsConfig)

	// TODO: paginate?
	resp, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.BucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to list items in bucket %v", err)
	}

	for _, item := range resp.Contents {
		fmt.Printf("Name: %s, Size: %d bytes\n", *item.Key, item.Size)
	}

	return nil
}
