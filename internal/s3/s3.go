package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/supanova-rp/supanova-file-cleaner/internal/config"
)

type Client struct {
	s3         *s3.Client
	bucketName string
}

type Item struct {
	Key  string
	Size int64
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
		s3:         client,
		bucketName: cfg.BucketName,
	}, nil
}

func (c *Client) GetBucketItems(ctx context.Context) ([]Item, error) {
	paginator := s3.NewListObjectsV2Paginator(c.s3, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
	})

	items := make([]Item, 0)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to list items in bucket: %w", err)
		}

		for _, item := range page.Contents {
			items = append(items, Item{
				Key:  aws.ToString(item.Key),
				Size: aws.ToInt64(item.Size),
			})
		}
	}

	return items, nil
}

func (c *Client) DeleteItem(ctx context.Context, key string) error {
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete item with key: %s, error: %v", key, err)
	}

	return nil
}
