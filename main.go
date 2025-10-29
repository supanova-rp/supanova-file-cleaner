package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"

	"github.com/supanova-rp/supanova-file-cleaner/internal/config"
	"github.com/supanova-rp/supanova-file-cleaner/internal/filecleaner"
	"github.com/supanova-rp/supanova-file-cleaner/internal/s3"
	"github.com/supanova-rp/supanova-file-cleaner/internal/store"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("run failed:", err)
		os.Exit(1)
	}

	slog.Info("shutting down gracefully...")
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	cleaner := filecleaner.New(db, s3Client, cfg.DryRun)

	c := cron.New()

	_, err = c.AddFunc(cfg.CronSchedule, func() {
		// Use a separate context for the job so it isn't cancelled during shutdown and
		// will allow c.Stop() to wait for it to finish
		jobCtx := context.Background()
		err = cleaner.Run(jobCtx)
		if err != nil {
			slog.Error("file cleaner run failed", slog.Any("err", err))
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron func: %v", err)
	}

	c.Start()
	slog.Info("cron scheduler started", slog.String("schedule", cfg.CronSchedule))

	<-ctx.Done() // Blocks until signal received (e.g. by ctrl-C or process killed)

	cronCtx := c.Stop() // Returns a context that waits for any running jobs to finish, then sends to the ctx Done channel
	<-cronCtx.Done()

	slog.Info("cron scheduler stopped")
	return nil
}
