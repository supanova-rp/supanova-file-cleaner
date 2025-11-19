package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	CronSchedule string
	DryRun       bool // if DryRun is true, unused s3 items aren't actually deleted
	AWS          AWS
}

type AWS struct {
	AssetsBucketName   string
	DBBackupBucketName string
	Config             AWSConfig
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

func ParseEnv() (*Config, error) {
	// Ignore error because in production there will be no .env file, env vars will be passed
	// in at runtime via docker run command
	_ = godotenv.Load()

	envVars := map[string]string{
		"DATABASE_URL":                    "",
		"CRON_SCHEDULE":                   "",
		"AWS_REGION":                      "",
		"AWS_SUPANOVA_ASSETS_BUCKET_NAME": "",
		"AWS_DB_BACKUP_BUCKET_NAME":       "",
		"AWS_ACCESS_KEY_ID":               "",
		"AWS_SECRET_ACCESS_KEY":           "",
		"DRY_RUN":                         "",
	}

	for key := range envVars {
		value := os.Getenv(key)
		if value == "" {
			return nil, fmt.Errorf("%s environment variable is not set", key)
		}
		envVars[key] = value
	}

	dryRun, err := strconv.ParseBool(envVars["DRY_RUN"])
	if err != nil {
		return nil, fmt.Errorf("unable to parse DRY_RUN environment variable: %v", err)
	}

	return &Config{
		DatabaseURL:  envVars["DATABASE_URL"],
		CronSchedule: envVars["CRON_SCHEDULE"],
		DryRun:       dryRun,
		AWS: AWS{
			AssetsBucketName:   envVars["AWS_SUPANOVA_ASSETS_BUCKET_NAME"],
			DBBackupBucketName: envVars["AWS_DB_BACKUP_BUCKET_NAME"],
			Config: AWSConfig{
				Region:          envVars["AWS_REGION"],
				AccessKeyID:     envVars["AWS_ACCESS_KEY_ID"],
				SecretAccessKey: envVars["AWS_SECRET_ACCESS_KEY"],
			},
		},
	}, nil
}
