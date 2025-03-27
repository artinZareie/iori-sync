package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/artinZareie/iori-sync/internal/filesystem"
)

func TestWalkAsList(t *testing.T) {
	dir, err := os.MkdirTemp("", "walkaslist_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	testFiles := []string{
		"file1.txt",
		"file2.txt",
		filepath.Join("subdir", "file3.txt"),
	}

	for _, name := range testFiles {
		fullPath := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := filesystem.WalkAsList(dir)
	if err != nil {
		t.Fatal(err)
	}

	expectedPaths := []string{
		dir,
		filepath.Join(dir, "file1.txt"),
		filepath.Join(dir, "file2.txt"),
		filepath.Join(dir, "subdir"),
		filepath.Join(dir, "subdir", "file3.txt"),
	}

	if len(files) != len(expectedPaths) {
		t.Errorf("expected %d entries, got %d", len(expectedPaths), len(files))
		for _, file := range files {
			t.Logf("Path: %s, Name: %s, IsDir: %v", file.Path, file.Info.Name(), file.Info.IsDir())
		}
	}

	expectedSet := make(map[string]bool, len(expectedPaths))
	for _, path := range expectedPaths {
		expectedSet[path] = false
	}

	for _, file := range files {
		if _, exists := expectedSet[file.Path]; exists {
			expectedSet[file.Path] = true
		} else {
			t.Errorf("unexpected path: %s", file.Path)
		}
	}

	for path, found := range expectedSet {
		if !found {
			t.Errorf("missing expected path: %s", path)
		}
	}
}
