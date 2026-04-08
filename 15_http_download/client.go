package http_download

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads a URL to a local file path.
// It should:
//   - Support resuming partial downloads using HTTP Range requests.
//   - Verify the downloaded file's SHA-256 matches expectedHash (hex string).
//   - Return an error if the server returns a non-2xx status code.
//   - Clean up partial files on checksum failure.

func DownloadFile(url, destPath, expectedHash string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		f.Close()
		return fmt.Errorf("download: %w", err)
	}

	// Verify checksum
	hash, err := fileHash(destPath)
	if err != nil {
		return fmt.Errorf("hash: %w", err)
	}

	if hash != expectedHash {
		return fmt.Errorf("checksum mismatch: got %s, want %s", hash, expectedHash)
	}

	return f.Close()
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
