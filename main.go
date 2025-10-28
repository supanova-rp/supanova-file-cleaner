package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/robfig/cron/v3"
	"github.com/supanova-rp/supanova-file-cleaner/internal/config"
	"github.com/supanova-rp/supanova-file-cleaner/internal/filecleaner"
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

	cleaner := filecleaner.New(db, s3Client)

	c := cron.New()

	_, err = c.AddFunc(cfg.CronSchedule, func() {
		err = cleaner.Run(ctx)
		if err != nil {
			slog.Error("file cleaner run failed", slog.Any("err", err))
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron func: %v", err)
	}

	// Start the scheduler (non-blocking)
	c.Start()

	// TODO: Best way to make the app block?
	select {}
}
