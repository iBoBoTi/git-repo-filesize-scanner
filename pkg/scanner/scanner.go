package scanner

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io/fs"
	"path/filepath"
	"sync"
)

const mb = 1024 * 1024

// FileExceedingThreshold records a file exceeding the threshold.
type FileExceedingThreshold struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	SizeMB string `json:"-"`
}

// Scan walks repository directory and returns files whose size exceeds thresholdBytes.
func Scan(ctx context.Context, logger *zap.Logger, repoDir string, thresholdBytes int64) ([]FileExceedingThreshold, error) {
	var fileExceedingThresholdMutex sync.Mutex
	fileExceedingThresholds := make([]FileExceedingThreshold, 0)

	walkFn := func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dirEntry.IsDir() {
			if dirEntry.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}
		info, err := dirEntry.Info()
		if err != nil {
			logger.Debug("skipping path", zap.String("path", path), zap.Error(err))
			return nil
		}
		sizeInBytes := info.Size()
		if sizeInBytes > thresholdBytes {
			rel, _ := filepath.Rel(repoDir, path)
			sizeMB := float64(sizeInBytes) / float64(mb)
			fileExceedingThresholdMutex.Lock()
			fileExceedingThresholds = append(fileExceedingThresholds, FileExceedingThreshold{Name: rel, Size: sizeInBytes, SizeMB: fmt.Sprintf("%.2f MB", sizeMB)})
			fileExceedingThresholdMutex.Unlock()
		}
		return nil
	}

	if err := filepath.WalkDir(repoDir, walkFn); err != nil {
		return nil, err
	}
	return fileExceedingThresholds, nil
}
