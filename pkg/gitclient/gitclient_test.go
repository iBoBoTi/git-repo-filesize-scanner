//go:build !race
// +build !race

package gitclient

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/stretchr/testify/require"
)

// createTestRepo initialises a test repo I can clone from.
func createTestRepo(t *testing.T) string {
	dir := t.TempDir()
	repo, err := git.PlainInit(dir, false)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	readme := filepath.Join(dir, "README.md")
	require.NoError(t, os.WriteFile(readme, []byte("hello"), 0o644))
	_, err = worktree.Add("README.md")
	require.NoError(t, err)

	_, err = worktree.Commit("init", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test-name",
			Email: "test-name@email.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
	return dir
}

func TestCloneRepo_OK(t *testing.T) {
	src := createTestRepo(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dst, err := CloneRepo(ctx, "file://"+src, "")
	require.NoError(t, err)
	require.DirExists(t, dst)
	require.FileExists(t, filepath.Join(dst, "README.md"))
}

func TestCloneRepo_FailsAfterRetries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := CloneRepo(ctx, "https://invalid.repo.com/repo-does-not-exist.git", "")
	require.Error(t, err)
}
