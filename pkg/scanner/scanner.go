package scanner

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const mb = 1024 * 1024

var MaxWorkers = runtime.NumCPU() * 2

// FileExceedingThreshold records a file exceeding the threshold.
type FileExceedingThreshold struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	SizeMB string `json:"-"`
}

// Scan walks repository directory and returns files whose size exceeds thresholdBytes.
func Scan(ctx context.Context, logger *zap.Logger, repoDir string, thresholdBytes int64) ([]FileExceedingThreshold, error) {
	jobs := make(chan string, 1024)
	fileExceedingThresholdsChan := make(chan FileExceedingThreshold, 256)

	var wg sync.WaitGroup
	for i := 0; i < MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				if ctx.Err() != nil {
					return
				}
				info, err := os.Stat(path)
				if err != nil {
					logger.Debug("skipping path", zap.String("path", path), zap.Error(err))
					continue
				}
				sizeInBytes := info.Size()
				if sizeInBytes > thresholdBytes {
					rel, _ := filepath.Rel(repoDir, path)
					fileExceedingThresholdsChan <- FileExceedingThreshold{
						Name:   rel,
						Size:   sizeInBytes,
						SizeMB: fmt.Sprintf("%.2f MB", float64(sizeInBytes)/float64(mb)),
					}
				}
			}
		}()
	}

	go func() { wg.Wait(); close(fileExceedingThresholdsChan) }()

	walkErr := filepath.WalkDir(repoDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}

		select {
		case jobs <- path:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	close(jobs)

	filesExceedingThresholds := make([]FileExceedingThreshold, 0, 128)
	for fileExceedingThresholds := range fileExceedingThresholdsChan {
		filesExceedingThresholds = append(filesExceedingThresholds, fileExceedingThresholds)
	}

	return filesExceedingThresholds, walkErr
}
