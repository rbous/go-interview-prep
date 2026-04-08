# AI Grading Instructions

You are grading a candidate's solutions to debugging exercises in Go and C.
The candidate was given buggy source files and asked to fix them.
**Test files (`_test.go`, `test_*.c`) should NOT be modified** — only the source files.

## General Grading Process

### Go exercises (01-10, 15-18)
1. For each exercise, read the candidate's modified `.go` file.
2. Run `go test -race -v ./XX_exercise_name/` and record pass/fail.
3. Run `go vet ./XX_exercise_name/` and check for warnings.
4. Evaluate the solution against the rubric below.
5. Assign a score per exercise (0-10) and provide brief feedback.

### C exercises (11-14)
1. For each exercise, read the candidate's modified `.c` file.
2. Run `cd XX_dir && make test` and record pass/fail.
3. Check for sanitizer warnings (ASAN for 11-13, TSAN for 14).
4. Evaluate the solution against the rubric below.
5. Assign a score per exercise (0-10) and provide brief feedback.

## Scoring Scale

- **10**: All tests pass with `-race`, clean `go vet`, idiomatic solution.
- **8-9**: All tests pass, minor style issues (e.g., unnecessary complexity).
- **5-7**: Most tests pass or correct approach with minor implementation errors.
- **3-4**: Shows understanding of the bug but fix is incomplete or introduces new issues.
- **1-2**: Attempted but fundamentally wrong approach.
- **0**: No changes or test files were modified.

## Per-Exercise Rubrics

### 01 — Race Condition (`01_race_condition/counter.go`)

**Bugs to fix:**
- Data race on `Counter.count` — multiple goroutines read/write without synchronization.
- `IncrementConcurrently` returns before goroutines finish (no wait mechanism).

**Acceptable fixes:**
- `sync.Mutex` protecting `count` in `Increment()` and `Value()` — **preferred**.
- `sync/atomic` for the counter — also acceptable.
- Must add `sync.WaitGroup` or similar to wait for all goroutines.

**Deductions:**
- -3 if race detector still reports issues.
- -2 if goroutines aren't waited on (even if tests sometimes pass).
- -1 if uses a global lock instead of per-Counter lock.

---

### 02 — Deadlock (`02_deadlock/transfer.go`)

**Bug to fix:**
- Classic lock-ordering deadlock: `Transfer(a, b)` locks a then b, while `Transfer(b, a)` locks b then a.

**Acceptable fixes:**
- Lock by account ID order (lower ID first) — **preferred/canonical**.
- Single global mutex — correct but not ideal (-2 for lost concurrency).
- `TryLock` with retry — acceptable if correct.

**Deductions:**
- -5 if deadlock still possible.
- -2 if uses a single global lock (correct but poor concurrency).
- -1 if doesn't handle the `from == to` edge case (though tests don't check this).

---

### 03 — Goroutine Leak (`03_goroutine_leak/fetcher.go`)

**Bugs to fix:**
- Failed fetches don't send to `ch`, so the `for range urls` loop blocks waiting for results that will never come.
- Even successful goroutines may block on `ch <- result` if the function returns early via timeout.

**Acceptable fixes:**
- Use a buffered channel of size `len(urls)` so sends never block — **simplest**.
- Use `context` or `done` channel to signal goroutines to abort.
- Count expected successes differently (e.g., use a `sync.WaitGroup` + close pattern).

**Deductions:**
- -5 if goroutine leak test still fails.
- -2 if only fixes the error case but not the timeout case (or vice versa).
- -1 if solution is overly complex for what's needed.

---

### 04 — Channel Panic (`04_channel_panic/dispatcher.go`)

**Bugs to fix in `Dispatch`:**
- Multiple workers each call `close(resultCh)` — only one close is valid, rest panic.

**Fix:** Use a `sync.WaitGroup` to close `resultCh` once after all workers finish.

**Bugs to fix in `DispatchOrdered`:**
- `resultCh` is never closed, so `for r := range resultCh` blocks forever.

**Fix:** Add a goroutine that waits for all workers then closes `resultCh`.

**Deductions:**
- -3 per function that still panics or deadlocks.
- -2 if `DispatchOrdered` doesn't preserve order.
- -1 if channel close is done in a racy way.

---

### 05 — WaitGroup Bug (`05_waitgroup_bug/pipeline.go`)

**Bugs to fix:**
- `wg.Add(1)` is inside the goroutine — `wg.Wait()` may return before goroutines call `Add`.
- `results = append(results, ...)` is a data race (concurrent slice append).

**Acceptable fixes:**
- Move `wg.Add(1)` before `go func()`.
- Protect `results` with a `sync.Mutex`, OR use a channel to collect results, OR use indexed writes to a pre-allocated slice.

**Deductions:**
- -3 if `wg.Add` is still inside the goroutine.
- -3 if slice append is still racy.
- -1 if uses overly complex synchronization.

---

### 06 — Concurrent Map (`06_concurrent_map/registry.go`)

**Bugs to fix:**
- All map operations are unprotected — concurrent read/write panics.
- `List()` returns the internal map directly — callers can mutate internal state.

**Acceptable fixes:**
- `sync.RWMutex` with `RLock` for reads and `Lock` for writes — **preferred**.
- `sync.Map` — acceptable but less idiomatic for this use case.
- `List()` must return a copy of the map.

**Deductions:**
- -4 if concurrent access still panics.
- -2 if `List()` still returns the internal map reference.
- -1 if uses `Mutex` instead of `RWMutex` (correct but suboptimal for read-heavy workloads).

---

### 07 — Context Cancellation (`07_context_cancel/downloader.go`)

**Bugs to fix:**
- `simulateDownload` ignores context — sleeps for full duration even after cancel.
- Results slice has a data race (concurrent append without sync).

**Acceptable fixes:**
- Pass `context.Context` to `simulateDownload` and use `select` with `ctx.Done()` and `time.After`.
- Use mutex or channel for collecting results.
- `DownloadPackages` should return when context is cancelled (not wait for all goroutines).

**Deductions:**
- -4 if context cancellation is not respected (timing test fails).
- -3 if data race on results.
- -1 if goroutines leak after context cancel (nice to handle but not strictly tested).

---

### 08 — Atomic File Write (`08_atomic_file_write/updater.go`)

**Bugs to fix:**
- `WriteConfig` writes directly to the target file — not atomic.
- `EnsureDir` uses `0777` permissions — too open.

**Acceptable fixes for `WriteConfig`:**
- Write to a temp file in the same directory, then `os.Rename` — **canonical pattern**.
- Must handle cleanup of temp file on write errors.

**Acceptable fixes for `EnsureDir`:**
- Change `0777` to `0755`.

**Deductions:**
- -4 if write is still non-atomic (no temp+rename).
- -2 if temp file leaks on error.
- -1 if directory permissions still 0777.
- +1 bonus if they sync/fsync the temp file before rename.

---

### 09 — Slice Gotcha (`09_slice_gotcha/versions.go`)

**Bug to fix:**
- `versions[:0]` reuses the input slice's backing array — appending to it overwrites the caller's data.

**Acceptable fixes:**
- Allocate a new slice: `var result []string` or `result := make([]string, 0)` — **simplest**.
- Use `slices.Clone` or explicit copy before slicing.

**Both `FilterVersions` and `UniqueVersions` have the same bug.**

**Deductions:**
- -3 per function that still mutates the original slice.
- -1 if allocates unnecessarily large slice (e.g., `make([]string, len(versions))`).

---

### 10 — Graceful Shutdown (`10_graceful_shutdown/server.go`)

**Bugs to fix:**
- `Shutdown` doesn't wait for workers to finish processing jobs in the channel.
- Data race on `results` slice (workers append concurrently).
- `Submit` doesn't check if server is shutting down — sends to closed channel panics.

**Acceptable fixes:**
- Add `sync.WaitGroup` tracking for workers; `Shutdown` waits on it (with context deadline).
- Protect `results` with a mutex or use a results channel.
- Add a `stopped` flag (atomic bool or mutex-protected) checked in `Submit`.
- `Submit` should recover from send-on-closed-channel panic OR check the flag before sending.

**Deductions:**
- -3 if `Shutdown` doesn't wait for in-flight jobs.
- -3 if data race on results.
- -2 if `Submit` panics after shutdown.
- -1 if context timeout in `Shutdown` is not respected.

---

---

### 11 — C Buffer Overflow (`11_c_buffer_overflow/version_parser.c`)

**Bugs to fix:**
- `strcpy` into fixed-size buffers without length checking — buffer overflow.
- `compare_versions` doesn't handle NULL input — segfault.
- `format_version` uses `sprintf` without checking `buf_size` — overflow.

**Acceptable fixes:**
- Use `strncpy` or `snprintf` with size limits.
- Check string length before copying; return -1 if field exceeds `VERSION_FIELD_MAX`.
- NULL checks at top of `compare_versions`.
- Use `snprintf(buf, buf_size, ...)` and check return value against `buf_size`.

**Deductions:**
- -3 if ASAN still reports overflow.
- -2 per function with missing NULL/bounds check.
- -1 if uses `strncpy` without null-terminating.

---

### 12 — C Memory Leak (`12_c_memory_leak/manifest.c`)

**Bugs to fix:**
- Error path in `parse_manifest` doesn't free entries already added to the list.
- `free_manifest` doesn't free `entry->package_name` and `entry->version`.
- `find_entry` doesn't check for NULL manifest — segfault.

**Acceptable fixes:**
- On parse error, walk the linked list and free all entries + their strings before returning NULL.
- Add `free(curr->package_name)` and `free(curr->version)` in `free_manifest`.
- NULL check at top of `find_entry`.

**Deductions:**
- -4 if ASAN reports leaks.
- -2 if error path still leaks.
- -2 if `free_manifest` still leaks strings.
- -1 if `find_entry` crashes on NULL.

---

### 13 — C Use-After-Free (`13_c_use_after_free/update_queue.c`)

**Bugs to fix:**
- `queue_drain_finished` accesses `curr->next` after `queue_remove` frees `curr`.
- `queue_remove` dereferences `node->prev` and `node->next` without NULL checks (crashes on head/tail removal).
- `queue_destroy` doesn't free `node->package_name` and `node->version`.

**Acceptable fixes:**
- Save `curr->next` before calling `queue_remove`.
- In `queue_remove`, check if node is head (update `q->head`) and if node is tail (update `q->tail`).
- Free string fields in `queue_destroy`.

**Deductions:**
- -4 if ASAN reports use-after-free.
- -3 if head/tail removal crashes.
- -2 if `queue_destroy` still leaks.
- -1 if edge cases (single element, empty queue) aren't handled.

---

### 14 — C Pthread Race (`14_c_pthread_race/progress.c`)

**Bugs to fix:**
- `num_complete++` from multiple threads without synchronization.
- `bytes_downloaded` read/written from different threads without locking.
- Loop variable `i` captured by pointer — all threads may read the same index.
- `args` array is stack-allocated in the loop and reused — threads read stale data.

**Acceptable fixes:**
- Use `pthread_mutex_t` or `__atomic` builtins for `num_complete`.
- Mutex or atomics for `bytes_downloaded` (or accept per-entry access since each thread writes its own entry — explain this reasoning).
- Heap-allocate a struct per thread containing `{tracker, index}`.

**Deductions:**
- -3 if TSAN reports races.
- -3 if loop variable capture bug not fixed.
- -2 if `args` array lifetime issue not fixed.
- -1 if uses overly coarse locking (single mutex for all entries).

---

### 15 — HTTP Download (`15_http_download/client.go`)

**Bugs to fix:**
- No Range request for resume — always overwrites existing partial file.
- No HTTP status code check — 4xx/5xx responses silently written to disk.
- File not closed before computing checksum (data may not be flushed).
- Corrupt file not removed on checksum failure.

**Acceptable fixes:**
- Check if `destPath` exists; if so, get its size and set `Range: bytes=N-` header.
- Open file with `O_APPEND` or `O_WRONLY` at offset for resume.
- Check `resp.StatusCode` is 200 or 206; error otherwise.
- `f.Close()` before `fileHash()`.
- `os.Remove(destPath)` on checksum mismatch.

**Deductions:**
- -3 if resume doesn't work.
- -2 if bad HTTP status not caught.
- -2 if corrupt file left on disk.
- -1 if file not closed before hash.

---

### 16 — TCP Server (`16_tcp_server/server.go`)

**Bugs to fix:**
- Connections not closed when client disconnects (resource leak).
- `Shutdown` doesn't close listener — `Accept` blocks forever.
- No connection read deadline — slow clients hold goroutines.
- `connCount` modified without synchronization.

**Acceptable fixes:**
- `defer conn.Close()` in `handleConn`.
- `s.listener.Close()` in `Shutdown`.
- `conn.SetDeadline()` or `conn.SetReadDeadline()` in `handleConn`.
- Use `atomic.AddInt32` or mutex for `connCount`.

**Deductions:**
- -3 if `Shutdown` doesn't close listener.
- -2 if connections leak.
- -2 if `connCount` has data race.
- -1 if no deadline on connections.

---

### 17 — Process Execution (`17_process_exec/runner.go`)

**Bugs to fix:**
- `RunCommand` doesn't enforce timeout — uses `cmd.Run()` without context.
- Error return doesn't include stderr content.
- `RunScript` ignores timeout parameter.
- `RunWithRetry` doesn't check context cancellation between retries.

**Acceptable fixes:**
- Use `exec.CommandContext(ctx, ...)` with `context.WithTimeout`.
- Wrap error: `fmt.Errorf("command failed: %s: %w", errBuf.String(), err)`.
- `RunScript` should also use `CommandContext`.
- Check `ctx.Err()` at the top of each retry iteration.

**Deductions:**
- -3 if timeout not enforced (timing test fails).
- -2 if stderr not in error message.
- -2 if `RunWithRetry` doesn't respect context.
- -1 if `RunScript` timeout still ignored.

---

### 18 — File Locking (`18_file_locking/lock.go`)

**Bugs to fix:**
- `AcquireLock` uses blocking `LOCK_EX` — hangs if lock held.
- `Release` doesn't remove the lock file.

**Acceptable fixes:**
- Add `LOCK_NB` flag: `syscall.LOCK_EX | syscall.LOCK_NB`.
- `os.Remove(l.path)` in `Release` before or after closing.

**Deductions:**
- -4 if acquire still blocks.
- -3 if lock file not cleaned up.
- -1 if TOCTOU issue in `IsLocked` not documented (bonus if candidate adds a comment).
- +1 bonus if candidate checks that the fd still refers to the same inode.

---

## Final Report Format

```
## Exercise Results

| # | Exercise            | Tests Pass | Race Clean | Score | Notes           |
|---|---------------------|------------|------------|-------|-----------------|
| 01| Race Condition      | PASS/FAIL  | YES/NO     | X/10  | brief feedback  |
| 02| Deadlock            | PASS/FAIL  | YES/NO     | X/10  | brief feedback  |
| ...                                                                        |

## Summary
- Total: XX/100
- Strengths: ...
- Areas to improve: ...
- Overall assessment: [Strong / Adequate / Needs Work]

## Detailed Feedback
(Per-exercise explanations of what was done well and what could be improved.)
```

## Assessment Thresholds

Total is now out of 180 (18 exercises x 10 points). Normalize to percentage.

- **90-100%**: Excellent — strong systems programming skills, ready for the role.
- **75-89%**: Good — solid understanding with minor gaps.
- **60-74%**: Adequate — understands core concepts but needs more practice.
- **40-59%**: Below expectations — significant gaps.
- **Below 40%**: Not ready — fundamentals need work.
