package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	AWS         AWSConfig
}

type AWSConfig struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
}

func ParseEnv() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL environment variable is not set")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		return Config{}, errors.New("AWS_REGION environment variable is not set")
	}

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		return Config{}, errors.New("AWS_BUCKET_NAME environment variable is not set")
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return Config{}, errors.New("AWS_ACCESS_KEY_ID environment variable is not set")
	}

	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		return Config{}, errors.New("AWS_SECRET_ACCESS_KEY environment variable is not set")
	}

	return Config{
		DatabaseURL: databaseURL,
		AWS: AWSConfig{
			Region:          region,
			BucketName:      bucketName,
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		},
	}, nil
}
