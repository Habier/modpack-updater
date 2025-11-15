package updater

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	FoldersToRestore = []string{"config", "defaultconfigs", "modernfix", "mods", "schematics"}
)

// UpdateFolders updates the specified folders from a zip archive
func UpdateFolders(destinationDir, zipPath string) error {
	if err := validateInputs(destinationDir, zipPath); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	if err := removeOldFolders(destinationDir); err != nil {
		return fmt.Errorf("failed to remove old folders: %v", err)
	}

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

	for _, f := range r.File {
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

		targetPath := filepath.Join(destinationDir, filepath.FromSlash(f.Name))

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %v", targetPath, err)
		}

		dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", targetPath, err)
		}

		srcFile, err := f.Open()
		if err != nil {
			dstFile.Close()
			return fmt.Errorf("failed to open zip entry %s: %v", f.Name, err)
		}

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
