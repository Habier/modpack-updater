package updater

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createTestZip crea un archivo ZIP de prueba con una estructura de carpetas y archivos
func createTestZip(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()

	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create test zip file: %v", err)
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	for path, content := range files {
		f, err := w.Create(path)
		if err != nil {
			t.Fatalf("failed to create file in zip: %v", err)
		}

		if content != "" {
			_, err = f.Write([]byte(content))
			if err != nil {
				t.Fatalf("failed to write to file in zip: %v", err)
			}
		}
	}
}

func TestValidateInputs(t *testing.T) {
	tempDir := t.TempDir()
	zipFile := filepath.Join(tempDir, "test.zip")
	f, err := os.Create(zipFile)
	if err != nil {
		t.Fatalf("Failed to create test zip file: %v", err)
	}
	f.Close()

	tests := []struct {
		name          string
		destination   string
		zipPath       string
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid inputs",
			destination: tempDir,
			zipPath:     zipFile,
			expectError: false,
		},
		{
			name:          "empty destination",
			destination:   "",
			zipPath:       zipFile,
			expectError:   true,
			errorContains: "destination directory cannot be empty",
		},
		{
			name:          "empty zip path",
			destination:   tempDir,
			zipPath:       "",
			expectError:   true,
			errorContains: "ZIP file path cannot be empty",
		},
		{
			name:          "non-existent destination",
			destination:   filepath.Join(tempDir, "nonexistent"),
			zipPath:       zipFile,
			expectError:   true,
			errorContains: "destination directory does not exist",
		},
		{
			name:          "non-existent zip file",
			destination:   tempDir,
			zipPath:       filepath.Join(tempDir, "nonexistent.zip"),
			expectError:   true,
			errorContains: "ZIP file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInputs(tt.destination, tt.zipPath)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errorContains != "" && err.Error() != tt.errorContains &&
					!strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestRemoveOldFolders(t *testing.T) {
	tempDir := t.TempDir()

	// Create test directories
	for _, folder := range FoldersToRestore {
		err := os.MkdirAll(filepath.Join(tempDir, folder), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	t.Run("remove existing folders", func(t *testing.T) {
		err := removeOldFolders(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify folders were removed
		for _, folder := range FoldersToRestore {
			_, err := os.Stat(filepath.Join(tempDir, folder))
			if !os.IsNotExist(err) {
				t.Errorf("expected folder %s to be removed, but it still exists", folder)
			}
		}
	})

	t.Run("handle non-existent folders", func(t *testing.T) {
		err := removeOldFolders(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Should not return an error if folders don't exist
	})
}

func TestExtractFolders(t *testing.T) {
	tempDir := t.TempDir()
	testZip := filepath.Join(tempDir, "test.zip")

	// Crear un archivo ZIP de prueba
	testFiles := map[string]string{
		"config/test.cfg":  "config content",
		"mods/test.mod":    "mod content",
		"ignored/file.txt": "this should be ignored",
	}

	createTestZip(t, testZip, testFiles)

	t.Run("extract specific folders", func(t *testing.T) {
		extractDir := filepath.Join(tempDir, "extracted")
		err := os.MkdirAll(extractDir, 0755)
		if err != nil {
			t.Fatalf("failed to create extract directory: %v", err)
		}

		err = extractFolders(testZip, extractDir)
		if err != nil {
			t.Fatalf("extractFolders failed: %v", err)
		}

		// Verificar que los archivos correctos fueron extraídos
		for path, content := range testFiles {
			if !strings.HasPrefix(path, "ignored") {
				filePath := filepath.Join(extractDir, path)
				data, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("failed to read extracted file %s: %v", path, err)
					continue
				}
				if string(data) != content {
					t.Errorf("file content mismatch for %s. Expected '%s', got '%s'", path, content, string(data))
				}
			} else {
				// Verificar que los archivos ignorados no fueron extraídos
				_, err := os.Stat(filepath.Join(extractDir, path))
				if !os.IsNotExist(err) {
					t.Errorf("file %s should not have been extracted", path)
				}
			}
		}
	})

	t.Run("handle non-existent zip file", func(t *testing.T) {
		extractDir := filepath.Join(tempDir, "nonexistent_zip")
		err := extractFolders("nonexistent.zip", extractDir)
		if err == nil {
			t.Error("expected error for non-existent zip file, got nil")
		}
	})
}

func TestUpdateFolders(t *testing.T) {
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test_update.zip")

	// Crear archivos de prueba
	testFiles := map[string]string{
		"config/test.cfg": "config content",
		"mods/test.mod":   "mod content",
	}

	createTestZip(t, zipPath, testFiles)

	t.Run("successful update", func(t *testing.T) {
		destDir := filepath.Join(tempDir, "minecraft")
		err := os.MkdirAll(destDir, 0755)
		if err != nil {
			t.Fatalf("failed to create destination directory: %v", err)
		}

		err = UpdateFolders(destDir, zipPath)
		if err != nil {
			t.Fatalf("UpdateFolders failed: %v", err)
		}

		// Verificar que los archivos fueron extraídos correctamente
		for path := range testFiles {
			_, err := os.Stat(filepath.Join(destDir, path))
			if os.IsNotExist(err) {
				t.Errorf("expected file %s to exist, but it doesn't", path)
			}
		}
	})

	t.Run("invalid inputs", func(t *testing.T) {
		err := UpdateFolders("", "test.zip")
		if err == nil {
			t.Error("expected error for empty destination, got nil")
		}

		err = UpdateFolders("/tmp", "")
		if err == nil {
			t.Error("expected error for empty zip path, got nil")
		}
	})
}
