package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	// FoldersToRestore defines which folders should be updated
	FoldersToRestore = []string{"config", "defaultconfigs", "modernfix", "mods", "schematics"}
)

// UpdateFolders updates the specified folders from a zip archive
func UpdateFolders(destinationDir, zipPath string) error {
	// Validate inputs
	if err := validateInputs(destinationDir, zipPath); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Note: Ownership restoration is not supported on Windows
	fmt.Println("⚠ Note: File ownership preservation is not supported on this platform")

	// Remove old folders
	if err := removeOldFolders(destinationDir); err != nil {
		return fmt.Errorf("failed to remove old folders: %v", err)
	}

	// Extract folders from zip
	if err := extractFolders(zipPath, destinationDir); err != nil {
		return fmt.Errorf("failed to extract folders: %v", err)
	}

	fmt.Println("✅ Process completed successfully.")
	return nil
}

func validateInputs(destinationDir, zipPath string) error {
	if destinationDir == "" {
		return fmt.Errorf("destination directory cannot be empty")
	}

	if zipPath == "" {
		return fmt.Errorf("ZIP file path cannot be empty")
	}

	// Check if destination directory exists
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		return fmt.Errorf("destination directory does not exist: %s", destinationDir)
	}

	// Check if zip file exists and is readable
	file, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("ZIP file does not exist or is not readable: %s", zipPath)
	}
	file.Close()

	return nil
}

func removeOldFolders(destinationDir string) error {
	fmt.Println("🧹 Removing old folders...")
	for _, folder := range FoldersToRestore {
		target := filepath.Join(destinationDir, folder)
		if _, err := os.Stat(target); err == nil {
			if err := os.RemoveAll(target); err != nil {
				return fmt.Errorf("failed to remove %s: %v", target, err)
			}
			fmt.Printf("✔ Removed: %s\n", target)
		} else {
			fmt.Printf("⚠ Not found: %s\n", target)
		}
	}
	fmt.Println()
	return nil
}

func extractFolders(zipPath, destinationDir string) error {
	fmt.Println("📦 Extracting folders from ZIP...")

	// Open the zip file
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	// Create a map of folders we want to extract
	foldersToExtract := make(map[string]bool)
	for _, folder := range FoldersToRestore {
		foldersToExtract[folder+"/"] = true
	}

	// Extract files
	for _, f := range r.File {
		// Check if this file is in one of our target folders
		extract := false
		for folder := range foldersToExtract {
			if len(f.Name) >= len(folder) && f.Name[0:len(folder)] == folder {
				extract = true
				break
			}
		}

		if !extract {
			continue
		}

		// Create the directory structure
		targetPath := filepath.Join(destinationDir, filepath.FromSlash(f.Name))

		if f.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetPath, err)
			}
			continue
		}

		// Create parent directories if they don't exist
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %v", targetPath, err)
		}

		// Create the file
		dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", targetPath, err)
		}

		srcFile, err := f.Open()
		if err != nil {
			dstFile.Close()
			return fmt.Errorf("failed to open zip entry %s: %v", f.Name, err)
		}

		// Copy file contents
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			srcFile.Close()
			dstFile.Close()
			return fmt.Errorf("failed to extract file %s: %v", f.Name, err)
		}

		srcFile.Close()
		dstFile.Close()
	}

	fmt.Println()
	return nil
}

// restoreOwnership is a no-op on Windows
func restoreOwnership() error {
	// This is a no-op on Windows
	return nil
}
