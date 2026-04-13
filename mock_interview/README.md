# Mock Interview Exercises

These exercises simulate the expected interview format:

1. **Design discussion** (5-10 min): Read the prompt at the top of each `.go` file and practice answering **out loud**.
2. **Bug fixing** (20-30 min): The interviewer gives you a buggy implementation. Fix it so all tests pass.

## Running Tests

```bash
# Single exercise
go test -race -v ./mock_interview/01_verify_package/

# All mock interviews
go test -race ./mock_interview/...
```

## Exercises

| # | Topic | Design Question | Bugs |
|---|-------|----------------|------|
| 01 | Package Verification | How do you verify an OTA update before installing? | 3 |
| 02 | Update Orchestration | How do you update components with dependency ordering? | 2 |
| 03 | Health Monitor & Rollback | After rebooting, how does the system commit or rollback? | 4 |

## Rules

- Do **NOT** modify the test files.
- Run with `-race` — some bugs only surface under the race detector.
- When you find a bug, **explain it out loud** before fixing it.
- There are usually more bugs than you think.
