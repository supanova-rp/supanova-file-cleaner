package main

import (
	"context"
	"fmt"
	"os"

	"github.com/supanova-rp/supanova-file-cleaner/internal/config"
	"github.com/supanova-rp/supanova-file-cleaner/internal/s3"
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

	s3Client, err := s3.New(ctx, cfg.AWS)
	if err != nil {
		return fmt.Errorf("unable to connect to s3: %v", err)
	}

	err = s3Client.ListBucket(ctx, cfg.AWS)
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
