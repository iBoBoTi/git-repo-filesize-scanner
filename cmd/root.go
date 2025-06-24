package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	rootCmd = &cobra.Command{
		Use:   "grfscan",
		Short: "Scan a GitHub repository for large files",
	}
	logger *zap.Logger
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if logger != nil {
			logger.Sugar().Fatal(err)
		}
		os.Exit(1)
	}
}

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(func() {
		cobra.EnableCommandSorting = true
	})
	rootCmd.AddCommand(scanCmd)
}
