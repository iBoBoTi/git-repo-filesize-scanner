package scanner

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestScan_FindsOnlyFilesAboveThreshold(t *testing.T) {
	old := MaxWorkers
	MaxWorkers = 2
	defer func() { MaxWorkers = old }()

	repo := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(repo, "small.txt"), make([]byte, 512), 0o644))

	large := filepath.Join(repo, "large.bin")
	require.NoError(t, os.WriteFile(large, make([]byte, mb+mb/2), 0o644))

	logger := zaptest.NewLogger(t)

	hits, err := Scan(context.Background(), logger, repo, int64(mb))
	require.NoError(t, err)

	require.Len(t, hits, 1, "only large.bin should be reported")
	require.Equal(t, "large.bin", hits[0].Name)
	require.Greater(t, hits[0].Size, int64(mb))
}

func TestScan_ContextCancelStopsEarly(t *testing.T) {
	old := MaxWorkers
	MaxWorkers = 4
	defer func() { MaxWorkers = old }()

	repo := t.TempDir()
	for i := 0; i < 5000; i++ {
		_ = os.WriteFile(filepath.Join(repo, "f"+strconv.Itoa(i)), []byte("x"), 0o644)
	}

	ctx, cancel := context.WithCancel(context.Background())
	logger := zaptest.NewLogger(t)

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	_, err := Scan(ctx, logger, repo, int64(1))
	require.ErrorIs(t, err, context.Canceled, "scan must return ctx error when canceled")
}
