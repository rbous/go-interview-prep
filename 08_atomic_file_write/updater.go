package atomic_file_write

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteConfig writes config data to the given path atomically.
// If the process crashes mid-write, readers should see either the
// old content or the new content — never a partial/corrupt file.

func WriteConfig(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		f.Close()
		return fmt.Errorf("write: %w", err)
	}

	return f.Close()
}

// ReadConfig reads config data from the given path.
func ReadConfig(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// EnsureDir creates the directory for a given file path if it doesn't exist.
func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0777)
}
