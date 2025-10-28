package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supanova-rp/supanova-file-cleaner/internal/config"
)

type Client struct {
	s3 *s3.Client
}

func New(ctx context.Context, cfg config.AWSConfig) (*Client, error) {
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
		return nil, fmt.Errorf("unable to load SDK config %v", err)
	}

	client := s3.NewFromConfig(awsConfig)

	return &Client{
		s3: client,
	}, nil
}

func (c *Client) ListBucket(ctx context.Context, cfg config.AWSConfig) error {
	// TODO: paginate?
	resp, err := c.s3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
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
