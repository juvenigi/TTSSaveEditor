package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReconcile(t *testing.T) {
	// Create temp src and dst dirs
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")
	if err2 := os.MkdirAll(srcDir, 0755); err2 != nil {
		t.Fatal(err2)
	}
	if err2 := os.MkdirAll(dstDir, 0755); err2 != nil {
		t.Fatal(err2)
	}

	// Setup files in src
	err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	srcFiles := []string{
		"file1.txt",
		filepath.Join("subdir", "file2.txt"),
	}
	for _, f := range srcFiles {
		if err := os.WriteFile(filepath.Join(srcDir, f), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	err = Reconcile(srcDir, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := DeepLs(dstDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(actual) != 2 {
		t.Fatalf("expected 2 files got: %v", dstDir)
	}

	t.Log("ready")
}

func TestDeepLsMaxDepthExceeded(t *testing.T) {
	dir := t.TempDir()

	// Create nested dirs deeper than maxDepth
	deepPath := dir
	for i := 0; i <= MaxDepth+1; i++ {
		deepPath = filepath.Join(deepPath, "subdir")
	}
	if err := os.MkdirAll(deepPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a file at deepest level
	err := os.WriteFile(filepath.Join(deepPath, "file.txt"), []byte("data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// deepLs should return max depth exceeded error
	_, err = DeepLs(dir)
	if err == nil || !strings.Contains(err.Error(), "maximum depth exceeded") {
		t.Errorf("expected max depth exceeded error but got %v", err)
	}
}
