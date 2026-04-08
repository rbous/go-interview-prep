# Go Interview Prep: Linux Update Systems (Tesla)

10 exercises focused on debugging Go code, with emphasis on concurrency bugs
and systems-level patterns relevant to Linux update systems.

## Setup

```bash
cd go-interview-prep
go mod tidy
```

## How to work through exercises

Each exercise contains a `.go` file with **buggy code** and a `_test.go` file with tests.
Your job is to **fix the bugs** so all tests pass.

**Do not modify the test files.** Only edit the main `.go` files.

### Running tests

```bash
# Run a single exercise
go test ./01_race_condition/

# Run with race detector (important for concurrency exercises)
go test -race ./01_race_condition/

# Run all exercises
go test -race ./...
```

## Exercises

| #  | Name                | Category       | Difficulty |
|----|---------------------|----------------|------------|
| 01 | Race Condition      | Concurrency    | Easy       |
| 02 | Deadlock            | Concurrency    | Medium     |
| 03 | Goroutine Leak      | Concurrency    | Medium     |
| 04 | Channel Panic       | Concurrency    | Easy       |
| 05 | WaitGroup Bug       | Concurrency    | Easy       |
| 06 | Concurrent Map      | Concurrency    | Medium     |
| 07 | Context Cancellation| Concurrency    | Medium     |
| 08 | Atomic File Write   | Systems        | Medium     |
| 09 | Slice Gotcha        | General Go     | Easy       |
| 10 | Graceful Shutdown   | Systems/Conc.  | Hard       |

## Tips

- Always run tests with `-race` for concurrency exercises.
- Read the test file to understand the expected behavior.
- Each bug file has comments describing what the function **should** do.
- Some exercises have multiple bugs.

## Grading

See `GRADING.md` for AI-grading instructions.
