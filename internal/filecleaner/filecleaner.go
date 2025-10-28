package filecleaner

import (
	"context"
	"fmt"
	"strings"

	"github.com/supanova-rp/supanova-file-cleaner/internal/s3"
	"github.com/supanova-rp/supanova-file-cleaner/internal/store"
	"github.com/supanova-rp/supanova-file-cleaner/internal/store/sqlc"
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

	for _, item := range items {
		if isVideoItem(item.Key) && !isInVideoTable(item.Key, videos) {
			unusedItems = append(unusedItems, item)
		}

		if isCourseMaterialItem(item.Key) && !isInCourseMaterialsTable(item.Key, courseMaterials) {
			unusedItems = append(unusedItems, item)
		}
	}

	// for _, item := range items {
	// 	fmt.Printf("Name: %s, Size: %d bytes\n", item.Key, item.Size)
	// }

	// fmt.Println(">>> len(items): ", len(items))

	// for _, video := range videos {
	// 	fmt.Printf("Video in DB >> Name: %s, CourseID: %s, StorageKey: %s\n", video.Title.String, video.CourseID.String(), video.StorageKey.String())
	// }

	// for _, m := range courseMaterials {
	// 	fmt.Printf("Course Material in DB >> Name: %s, CourseID: %s, StorageKey: %s\n", m.Name, m.CourseID.String(), m.StorageKey.String())
	// }

	var totalSize int64

	for _, item := range unusedItems {
		totalSize += item.Size
		fmt.Printf("Unused Item >> Name: %s, Size: %d bytes\n", item.Key, item.Size)
	}

	// TODO: change to slog
	fmt.Printf("Number of deleted items: %d\n", len(unusedItems))
	fmt.Printf("Total size of deleted items: %dMB\n", totalSize/1000000)

	// TODO: delete the items

	return nil
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
