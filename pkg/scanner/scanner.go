package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
)

const mb = 1024 * 1024

// FileExceedingThreshold records a file exceeding the threshold.
type FileExceedingThreshold struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

// Scan walks repository directory and returns files whose size exceeds thresholdBytes.
func Scan(ctx context.Context, repoDir string, thresholdBytes int64) ([]FileExceedingThreshold, error) {
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
			return err
		}
		if info.Size() > thresholdBytes {
			rel, _ := filepath.Rel(repoDir, path)
			sizeMB := float64(info.Size()) / float64(mb)
			fileExceedingThresholdMutex.Lock()
			fileExceedingThresholds = append(fileExceedingThresholds, FileExceedingThreshold{Name: rel, Size: fmt.Sprintf("%.2f MB", sizeMB)})
			fileExceedingThresholdMutex.Unlock()
		}
		return nil
	}

	if err := filepath.WalkDir(repoDir, walkFn); err != nil {
		return nil, err
	}
	return fileExceedingThresholds, nil
}
