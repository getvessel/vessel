package tests

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"testing"

	"vessel.dev/vessel/internal/utils"
)

func TestCreateTarContext(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "vessel_tar_test")
	_ = os.RemoveAll(tempDir)
	_ = os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "subdir")
	_ = os.MkdirAll(subDir, 0755)

	filePath := filepath.Join(tempDir, "hello.txt")
	_ = os.WriteFile(filePath, []byte("Hello, Vessel!"), 0644)

	subFilePath := filepath.Join(subDir, "nested.txt")
	_ = os.WriteFile(subFilePath, []byte("Nested file content"), 0644)

	reader, err := utils.CreateTarContext(tempDir)
	if err != nil {
		t.Fatalf("CreateTarContext failed: %v", err)
	}

	tr := tar.NewReader(reader)
	filesFound := make(map[string]string)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed to read from tar archive: %v", err)
		}

		if header.Typeflag == tar.TypeReg {
			content, err := io.ReadAll(tr)
			if err != nil {
				t.Fatalf("failed to read file content for %s: %v", header.Name, err)
			}
			filesFound[header.Name] = string(content)
		} else if header.Typeflag == tar.TypeDir {
			filesFound[header.Name] = "DIR"
		}
	}

	if content, ok := filesFound["hello.txt"]; !ok || content != "Hello, Vessel!" {
		t.Fatalf("unexpected or missing hello.txt: %q", content)
	}
	if content, ok := filesFound["subdir/nested.txt"]; !ok || content != "Nested file content" {
		t.Fatalf("unexpected or missing subdir/nested.txt: %q", content)
	}
}
