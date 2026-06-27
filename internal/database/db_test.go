package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPortableFlagInBuildDirectory(t *testing.T) {
	root := t.TempDir()
	exeDir := filepath.Join(root, "build", "bin", "postme.app", "Contents", "MacOS")
	mustMkdirAll(t, exeDir)
	mustWriteFile(t, filepath.Join(root, "build", portableFlagName), []byte("test-data"))

	dataDir, ok := getPortableDataDir(exeDir)
	if !ok {
		t.Fatal("expected portable mode")
	}

	want := filepath.Join(root, "build", "test-data")
	if dataDir != want {
		t.Fatalf("dataDir = %q, want %q", dataDir, want)
	}
}

func TestPortableFlagInBuildBinDirectory(t *testing.T) {
	root := t.TempDir()
	exeDir := filepath.Join(root, "build", "bin", "postme.app", "Contents", "MacOS")
	mustMkdirAll(t, exeDir)
	mustWriteFile(t, filepath.Join(root, "build", "bin", portableFlagName), nil)

	dataDir, ok := getPortableDataDir(exeDir)
	if !ok {
		t.Fatal("expected portable mode")
	}

	want := filepath.Join(root, "build", "bin", dataDirName)
	if dataDir != want {
		t.Fatalf("dataDir = %q, want %q", dataDir, want)
	}
}

func TestPortableFlagClosestDirectoryWins(t *testing.T) {
	root := t.TempDir()
	exeDir := filepath.Join(root, "build", "bin", "postme.app", "Contents", "MacOS")
	mustMkdirAll(t, exeDir)
	mustWriteFile(t, filepath.Join(root, "build", portableFlagName), []byte("outer-data"))
	mustWriteFile(t, filepath.Join(exeDir, portableFlagName), []byte("inner-data"))

	dataDir, ok := getPortableDataDir(exeDir)
	if !ok {
		t.Fatal("expected portable mode")
	}

	want := filepath.Join(exeDir, "inner-data")
	if dataDir != want {
		t.Fatalf("dataDir = %q, want %q", dataDir, want)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}
