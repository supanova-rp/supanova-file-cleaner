package backup

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/supanova-rp/supanova-maintenance/internal/s3"
)

type Backup struct {
	dbURL string
	s3    *s3.Client
}

func New(dbURL string, s3Client *s3.Client) *Backup {
	return &Backup{
		dbURL: dbURL,
		s3:    s3Client,
	}
}

func (b *Backup) Run(ctx context.Context) error {
	slog.Info("running db backup")

	config, err := pgx.ParseConfig(b.dbURL)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	cmd := exec.Command("pg_dump",
		"-h", config.Host,
		"-p", fmt.Sprintf("%d", config.Port),
		"-U", config.User,
		"-d", config.Database,
	)

	// Set password via environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", config.Password))

	// Execute pg_dump and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %w\nOutput: %s", err, string(output))
	}

	key := fmt.Sprintf("backup_%s.sql", time.Now().Format("02-01-2006"))

	err = b.s3.PutItem(ctx, key, output)
	if err != nil {
		return fmt.Errorf("db backup s3 upload failed: %v", err)
	}

	slog.Info("db backup successful")
	return nil
}
