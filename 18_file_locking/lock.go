package file_locking

import (
	"fmt"
	"os"
	"syscall"
)

// UpdateLock prevents multiple update processes from running simultaneously.
// Uses flock(2) for advisory file locking on Linux/macOS.
//
// Usage:
//   lock, err := AcquireLock("/var/run/updater.lock")
//   if err != nil { /* another updater is running */ }
//   defer lock.Release()

type UpdateLock struct {
	file *os.File
	path string
}

func AcquireLock(path string) (*UpdateLock, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create lock file: %w", err)
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, fmt.Errorf("flock: %w", err)
	}

	// Write PID for debugging
	fmt.Fprintf(f, "%d\n", os.Getpid())

	return &UpdateLock{file: f, path: path}, nil
}

func (l *UpdateLock) Release() error {
	if l.file == nil {
		return nil
	}

	if err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN); err != nil {
		return fmt.Errorf("unlock: %w", err)
	}

	return l.file.Close()
}

// IsLocked checks if the lock file at path is currently held by another process.

func IsLocked(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return true // couldn't get lock, someone else has it
	}

	// We got the lock, so it wasn't held. Release it.
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return false
}
