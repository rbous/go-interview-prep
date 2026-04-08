package http_download

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func makeTestServer(content string) *httptest.Server {
	data := []byte(content)
	hash := sha256.Sum256(data)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Support Range requests
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			// Parse "bytes=N-"
			parts := strings.TrimPrefix(rangeHeader, "bytes=")
			dashIdx := strings.Index(parts, "-")
			if dashIdx >= 0 {
				startStr := parts[:dashIdx]
				start, err := strconv.Atoi(startStr)
				if err == nil && start < len(data) {
					w.Header().Set("Content-Length", strconv.Itoa(len(data)-start))
					w.Header().Set("Content-Range",
						fmt.Sprintf("bytes %d-%d/%d", start, len(data)-1, len(data)))
					w.WriteHeader(http.StatusPartialContent)
					w.Write(data[start:])
					return
				}
			}
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("X-Checksum-Sha256", fmt.Sprintf("%x", hash))
		w.Write(data)
	}))
}

func TestDownloadBasic(t *testing.T) {
	content := "This is firmware version 2.1.0 payload data"
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	srv := makeTestServer(content)
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "firmware.bin")
	err := DownloadFile(srv.URL, dest, hash)
	if err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(dest)
	if string(got) != content {
		t.Errorf("content mismatch")
	}
}

func TestDownloadChecksumMismatch(t *testing.T) {
	content := "real firmware data"
	srv := makeTestServer(content)
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "firmware.bin")
	err := DownloadFile(srv.URL, dest, "badhash")
	if err == nil {
		t.Fatal("expected checksum error")
	}

	// Corrupt file should be cleaned up
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Error("corrupt file should be removed on checksum failure")
	}
}

func TestDownloadResume(t *testing.T) {
	content := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	srv := makeTestServer(content)
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "firmware.bin")

	// Write first 10 bytes as a "partial download"
	os.WriteFile(dest, []byte(content[:10]), 0644)

	err := DownloadFile(srv.URL, dest, hash)
	if err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(dest)
	if string(got) != content {
		t.Errorf("content after resume = %q, want %q", string(got), content)
	}
}

func TestDownloadBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "internal server error")
	}))
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "firmware.bin")
	err := DownloadFile(srv.URL, dest, "anyhash")
	if err == nil {
		t.Fatal("expected error on 500 status")
	}
}

func TestDownloadToNewFile(t *testing.T) {
	content := "fresh download"
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	srv := makeTestServer(content)
	defer srv.Close()

	dest := filepath.Join(t.TempDir(), "new_firmware.bin")
	err := DownloadFile(srv.URL, dest, hash)
	if err != nil {
		t.Fatal(err)
	}
}
