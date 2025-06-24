package domain

import (
	"encoding/json"
	"errors"
	"io"
)

type GitRepoLargeFileScanner struct {
	RepoCloneURL string  `json:"clone_url"`
	RepoSizeMB   float64 `json:"size"`
	Token        string  `json:"token,omitempty"`
}

// SizeBytes converts the repo file size threshold to bytes.
func (s GitRepoLargeFileScanner) SizeBytes() int64 {
	return int64(s.RepoSizeMB * 1024 * 1024)
}

// ParseJSON reads JSON input argument from r.
func ParseJSON(reader io.Reader) (GitRepoLargeFileScanner, error) {
	var repoLargerFileScanner GitRepoLargeFileScanner
	dec := json.NewDecoder(reader)
	if err := dec.Decode(&repoLargerFileScanner); err != nil {
		return GitRepoLargeFileScanner{}, err
	}

	if repoLargerFileScanner.RepoCloneURL == "" {
		return GitRepoLargeFileScanner{}, errors.New("clone_url is required")
	}
	if repoLargerFileScanner.RepoSizeMB <= 0 {
		return GitRepoLargeFileScanner{}, errors.New("size must be positive")
	}
	return repoLargerFileScanner, nil
}
