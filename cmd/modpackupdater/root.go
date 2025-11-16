package modpackupdater

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "modpack-updater",
	Short: "A tool to update modpack folders from a ZIP archive",
	Long: `A simple utility that removes a defined set of folders from a destination
directory and restores them from a ZIP backup.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
