package filecleaner

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/supanova-rp/supanova-file-cleaner/internal/s3"
	"github.com/supanova-rp/supanova-file-cleaner/internal/store"
	"github.com/supanova-rp/supanova-file-cleaner/internal/store/sqlc"
)

const (
	bytesInMB = 1000000
)

type FileCleaner struct {
	store *store.Store
	s3    *s3.Client
}

func New(db *store.Store, s3Client *s3.Client) *FileCleaner {
	return &FileCleaner{
		store: db,
		s3:    s3Client,
	}
}

func (f *FileCleaner) Run(ctx context.Context) error {
	slog.Info("Running file cleaner")

	items, err := f.s3.GetBucketItems(ctx)
	if err != nil {
		return fmt.Errorf("failed to list bucket: %v", err)
	}

	videos, err := f.store.Queries.GetVideos(ctx)
	if err != nil {
		return fmt.Errorf("failed to get videos: %v", err)
	}

	courseMaterials, err := f.store.Queries.GetCourseMaterials(ctx)
	if err != nil {
		return fmt.Errorf("failed to get course materials: %v", err)
	}

	var unusedItems []s3.Item
	var totalUnusedSize, totalUsedSize int64

	for _, item := range items {
		totalUsedSize += item.Size

		if isUnused(item, videos, courseMaterials) {
			unusedItems = append(unusedItems, item)
		}
	}

	slog.Info("Number of items in s3", slog.Int("count", len(items)))
	slog.Info("Total size of items in s3", slog.Int64("size_mb", totalUsedSize/bytesInMB))

	for _, item := range unusedItems {
		totalUnusedSize += item.Size
		slog.Info("Unused item", slog.String("name", item.Key), slog.Int64("size_bytes", item.Size))

		err := f.s3.DeleteItem(ctx, item.Key)
		if err != nil {
			return err
		}
	}

	if len(unusedItems) > 0 {
		slog.Info("Found unused items", slog.Int("count", len(unusedItems)))
		slog.Info("Total size of deleted items", slog.Int64("size_mb", totalUnusedSize/bytesInMB))
	} else {
		slog.Info("No unused items found")
	}

	return nil
}

func isUnused(item s3.Item, videos []sqlc.GetVideosRow, courseMaterials []sqlc.GetCourseMaterialsRow) bool {
	if isVideoItem(item.Key) && !isInVideoTable(item.Key, videos) {
		return true
	}

	if isCourseMaterialItem(item.Key) && !isInCourseMaterialsTable(item.Key, courseMaterials) {
		return true
	}

	return false
}

func isVideoItem(key string) bool {
	return strings.Contains(key, "/videos/")
}

func isCourseMaterialItem(key string) bool {
	// Only include PDFs to ensure if a new course material type is added, it won't be accidentally
	// deleted. This code will have to be updated to handle any new file types
	return strings.Contains(key, "/materials/") && strings.HasSuffix(key, ".pdf")
}

func isInVideoTable(key string, videos []sqlc.GetVideosRow) bool {
	courseID, videoID, found := strings.Cut(key, "/videos/")
	if !found {
		return false
	}

	for _, video := range videos {
		if courseID == video.CourseID.String() && videoID == video.StorageKey.String() {
			return true
		}
	}

	return false
}

func isInCourseMaterialsTable(key string, materials []sqlc.GetCourseMaterialsRow) bool {
	courseID, materialFileName, found := strings.Cut(key, "/materials/")
	if !found {
		return false
	}

	materialID := strings.TrimSuffix(materialFileName, ".pdf")

	for _, m := range materials {
		if courseID == m.CourseID.String() && materialID == m.StorageKey.String() {
			return true
		}
	}

	return false
}
