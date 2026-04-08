package process_exec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// RunCommand executes a shell command with a timeout and returns
// stdout, stderr, and any error. This is typical for update systems
// that need to run pre/post-install scripts.
//
// BUG(1): Timeout is not enforced — a hanging command blocks forever.
//         Should use context.WithTimeout.
// BUG(2): If the command fails (non-zero exit), the error doesn't include
//         stderr output, making debugging impossible.
// BUG(3): stdout and stderr buffers might not capture all output if the
//         command is killed by timeout (need to read pipes before Wait).

func RunCommand(command string, args []string, timeout time.Duration) (stdout string, stderr string, err error) {
	cmd := exec.Command(command, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	if err != nil {
		return outBuf.String(), errBuf.String(), err
	}

	return outBuf.String(), errBuf.String(), nil
}

// RunScript runs a shell script string via /bin/sh.
// Returns combined output and error.
//
// BUG: Does not use the timeout parameter at all.
// BUG: Does not capture stderr separately — hard to tell what went wrong.

func RunScript(script string, timeout time.Duration) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", script)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

// RunWithRetry runs a command up to `attempts` times, returning the first
// successful result. If all attempts fail, returns the last error.
//
// BUG: Ignores the context — should stop retrying if context is cancelled.
// BUG: No delay between retries — should use exponential backoff or at
//      least a fixed delay.

func RunWithRetry(ctx context.Context, command string, args []string, attempts int, timeout time.Duration) (string, error) {
	var lastErr error

	for i := 0; i < attempts; i++ {
		stdout, _, err := RunCommand(command, args, timeout)
		if err == nil {
			return stdout, nil
		}
		lastErr = err
	}

	return "", fmt.Errorf("all %d attempts failed: %w", attempts, lastErr)
}
