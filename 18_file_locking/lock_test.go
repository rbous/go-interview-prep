package file_locking

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAcquireAndRelease(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.lock")

	lock, err := AcquireLock(path)
	if err != nil {
		t.Fatal(err)
	}

	// Lock file should contain our PID
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("lock file should contain PID")
	}

	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestDoubleAcquireFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.lock")

	lock1, err := AcquireLock(path)
	if err != nil {
		t.Fatal(err)
	}
	defer lock1.Release()

	// Second acquire should fail immediately (non-blocking), not hang.
	done := make(chan error, 1)
	go func() {
		_, err := AcquireLock(path)
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Error("second AcquireLock should fail when lock is held")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("AcquireLock blocked — should return error immediately (use LOCK_NB)")
	}
}

func TestReleaseRemovesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.lock")

	lock, err := AcquireLock(path)
	if err != nil {
		t.Fatal(err)
	}

	lock.Release()

	// Lock file should be cleaned up
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("lock file should be removed after Release()")
	}
}

func TestAcquireAfterRelease(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.lock")

	lock1, err := AcquireLock(path)
	if err != nil {
		t.Fatal(err)
	}
	lock1.Release()

	// Should be able to acquire again
	lock2, err := AcquireLock(path)
	if err != nil {
		t.Fatalf("should be able to re-acquire after release: %v", err)
	}
	lock2.Release()
}

func TestReleaseNil(t *testing.T) {
	l := &UpdateLock{}
	if err := l.Release(); err != nil {
		t.Errorf("Release on zero-value should not error: %v", err)
	}
}

func TestIsLocked(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.lock")

	if IsLocked(path) {
		t.Error("should not be locked when file doesn't exist")
	}

	lock, err := AcquireLock(path)
	if err != nil {
		t.Fatal(err)
	}

	if !IsLocked(path) {
		t.Error("should report locked when lock is held")
	}

	lock.Release()
}
