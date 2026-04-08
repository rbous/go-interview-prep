# AI Grading Instructions

You are grading a candidate's solutions to Go debugging exercises.
The candidate was given buggy Go files and asked to fix them.
**Test files (`_test.go`) should NOT be modified** — only the main `.go` files.

## General Grading Process

1. For each exercise, read the candidate's modified `.go` file.
2. Run `go test -race -v ./XX_exercise_name/` and record pass/fail.
3. Run `go vet ./XX_exercise_name/` and check for warnings.
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

- **90-100**: Excellent — strong Go concurrency skills, ready for systems work.
- **75-89**: Good — solid understanding with minor gaps.
- **60-74**: Adequate — understands core concepts but needs more practice.
- **40-59**: Below expectations — significant gaps in concurrency understanding.
- **Below 40**: Not ready — fundamentals need work.
