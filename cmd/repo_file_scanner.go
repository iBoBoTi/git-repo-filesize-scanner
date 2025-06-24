package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/iBoBoTi/git-repo-filesize-scanner/internal/domain"
	git "github.com/iBoBoTi/git-repo-filesize-scanner/pkg/gitclient"
	"github.com/iBoBoTi/git-repo-filesize-scanner/pkg/scanner"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	fileInputPath string
	jsonArg       string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Clone a repo, scan files larger than the requested size threshold and return a summary about those files",
	RunE: func(cmd *cobra.Command, args []string) error {

		if jsonArg != "" && cmd.Flags().Changed("input") && fileInputPath != "-" {
			return fmt.Errorf("--json and --input cannot be used together")
		}

		var reader io.Reader
		switch {
		case jsonArg != "":
			reader = strings.NewReader(jsonArg)
		case fileInputPath == "-":
			reader = os.Stdin
		default:
			file, err := os.Open(filepath.Clean(fileInputPath))
			if err != nil {
				return fmt.Errorf("error opening file: %w", err)
			}
			defer file.Close()
			reader = file
		}

		gitRepoLargeFileScanner, err := domain.ParseJSON(reader)
		if err != nil {
			return fmt.Errorf("error parsing json: %w", err)
		}

		logger.Info("Cloning repository", zap.String("url", gitRepoLargeFileScanner.RepoCloneURL))
		repoDir, err := git.CloneRepo(cmd.Context(), gitRepoLargeFileScanner.RepoCloneURL, gitRepoLargeFileScanner.Token)
		if err != nil {
			return fmt.Errorf("error cloning repository: %w", err)
		}
		defer os.RemoveAll(repoDir)

		logger.Info("Scanning files", zap.String("dir", repoDir), zap.Float64("file_threshold_mb", gitRepoLargeFileScanner.RepoSizeMB))
		report, err := scanner.Scan(cmd.Context(), repoDir, gitRepoLargeFileScanner.SizeBytes())
		if err != nil {
			return fmt.Errorf("scan repo: %w", err)
		}

		output := struct {
			TotalNumOfFiles int                              `json:"total_num_of_files"`
			Files           []scanner.FileExceedingThreshold `json:"files"`
		}{
			TotalNumOfFiles: len(report),
			Files:           report,
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output)
	},
}

func init() {
	scanCmd.Flags().StringVarP(&fileInputPath, "input", "i", "-", "Path to JSON domain file or '-' for stdin (mutually exclusive with --json)")
	scanCmd.Flags().StringVarP(&jsonArg, "json", "j", "", "Inline JSON configuration string (mutually exclusive with --input)")
}
