//go:build test

package updater

import (
	"archive/zip"
	"os"
)

// createTestZip crea un archivo ZIP de prueba con una estructura de carpetas y archivos
func createTestZip(t test_helpers.TestingT, zipPath string, files map[string]string) {
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
