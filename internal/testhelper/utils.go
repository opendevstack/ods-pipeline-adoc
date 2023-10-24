package testhelper

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Golden(path string) string {
	return filepath.Join("../../test/testdata/golden", path)
}

func MkdirTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	return tempDir
}

func RmTempDir(tempDir string) {
	err := os.RemoveAll(tempDir)
	if err != nil {
		log.Printf("Failed to remove %s\n", tempDir)
	}
}

func CompareFiles(t *testing.T, wantFile, gotFile string) {
	t.Helper()
	got, err := os.ReadFile(gotFile)
	if err != nil {
		t.Error(err)
		return
	}
	want, err := os.ReadFile(wantFile)
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
