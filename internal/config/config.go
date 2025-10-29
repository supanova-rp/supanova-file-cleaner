package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	CronSchedule string
	DryRun       bool // if DryRun is true, unused s3 items aren't actually deleted
	AWS          AWSConfig
}

type AWSConfig struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
}

func ParseEnv() (*Config, error) {
	// Ignore error because in production there will be no .env file, env vars will be passed
	// in at runtime via docker run command
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is not set")
	}

	cronSchedule := os.Getenv("CRON_SCHEDULE")
	if cronSchedule == "" {
		return nil, errors.New("CRON_SCHEDULE environment variable is not set")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, errors.New("AWS_REGION environment variable is not set")
	}

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		return nil, errors.New("AWS_BUCKET_NAME environment variable is not set")
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil, errors.New("AWS_ACCESS_KEY_ID environment variable is not set")
	}

	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		return nil, errors.New("AWS_SECRET_ACCESS_KEY environment variable is not set")
	}

	dryRunString := os.Getenv("DRY_RUN")
	if dryRunString == "" {
		return nil, errors.New("DRY_RUN environment variable is not set")
	}
	dryRun, err := strconv.ParseBool(dryRunString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DRY_RUN environment variable: %v", err)
	}

	return &Config{
		DatabaseURL:  databaseURL,
		CronSchedule: cronSchedule,
		DryRun:       dryRun,
		AWS: AWSConfig{
			Region:          region,
			BucketName:      bucketName,
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		},
	}, nil
}
