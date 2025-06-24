package gitclient

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

// CloneRepo clones the given repository to a temporary directory using go-git.
func CloneRepo(ctx context.Context, cloneURL string, token string) (string, error) {
	dir, err := os.MkdirTemp("", "grfscan-*")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %w", err)
	}

	var auth transport.AuthMethod
	if token != "" {
		auth = &http.BasicAuth{
			Password: token,
		}
	}

	_, err = git.PlainCloneContext(ctx, dir, &git.CloneOptions{
		URL:      cloneURL,
		Depth:    1,
		Auth:     auth,
		Progress: os.Stdout,
	})
	if err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("error cloning repository: %w", err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return dir, nil
	}
	return absDir, nil
}
