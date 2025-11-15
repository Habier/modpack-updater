package main

import (
	"fmt"
	"path/filepath"

	"github.com/habier/modpack-updater/pkg/updater"
	"github.com/spf13/cobra"
)

var (
	destinationDir string
	zipPath        string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update modpack folders from a ZIP archive",
	Long: `Update modpack folders from a ZIP archive.

This command removes the specified folders from the destination directory
and extracts them from the provided ZIP archive.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		destDir, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("invalid destination directory: %v", err)
		}

		zipFile, err := filepath.Abs(args[1])
		if err != nil {
			return fmt.Errorf("invalid zip file path: %v", err)
		}

		// Print the operation details
		fmt.Printf("📂 Destination directory: %s\n", destDir)
		fmt.Printf("🗜 ZIP archive: %s\n\n", zipFile)

		// Call the updater
		return updater.UpdateFolders(destDir, zipFile)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you can define your flags and configuration settings
	updateCmd.Flags().StringVarP(&destinationDir, "destination", "d", "", "Destination directory (required)")
	updateCmd.Flags().StringVarP(&zipPath, "zip", "z", "", "Path to the ZIP archive (required)")
}
