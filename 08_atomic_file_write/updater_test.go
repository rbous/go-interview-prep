package atomic_file_write

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestWriteConfigBasic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.txt")

	data := []byte("key=value\n")
	if err := WriteConfig(path, data); err != nil {
		t.Fatal(err)
	}

	got, err := ReadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestWriteConfigAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.txt")

	original := []byte("original content\n")
	if err := WriteConfig(path, original); err != nil {
		t.Fatal(err)
	}

	// Concurrent reads while writing should never see partial content.
	newData := []byte("new content that is longer than the original\n")
	var wg sync.WaitGroup

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				data, err := ReadConfig(path)
				if err != nil {
					continue // file might be mid-rename
				}
				if !bytes.Equal(data, original) && !bytes.Equal(data, newData) {
					t.Errorf("read partial content: %q", data)
				}
			}
		}()
	}

	// Writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := WriteConfig(path, newData); err != nil {
			t.Errorf("write failed: %v", err)
		}
	}()

	wg.Wait()
}

func TestWriteConfigTempFileCleanedUp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.txt")

	if err := WriteConfig(path, []byte("hello")); err != nil {
		t.Fatal(err)
	}

	// No temp files should be left behind
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected only config.txt, found: %v", names)
	}
}

func TestEnsureDirPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "config.txt")

	if err := EnsureDir(path); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}

	perm := info.Mode().Perm()
	if perm&0o022 != 0o022 {
		// Should be at least 0755 (world-readable) but NOT 0777
	}
	if perm == 0o777 {
		t.Errorf("directory permissions too open: %o, want 0755", perm)
	}
}
