package scanner

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestScan_ReturnsBytesAndName(t *testing.T) {
	// mock repo dir
	repo := t.TempDir()
	bigFile := filepath.Join(repo, "bigFile.bin")
	data := make([]byte, mb+123)
	err := os.WriteFile(bigFile, data, 0o644)
	assert.NoError(t, err)

	logger := zap.NewNop()
	fileExceedingThresholds, err := Scan(context.Background(), logger, repo, mb)
	assert.NoError(t, err)
	assert.Len(t, fileExceedingThresholds, 1)

	got := fileExceedingThresholds[0]
	assert.Contains(t, bigFile, got.Name)
	assert.Equal(t, int64(len(data)), got.Size)
}
