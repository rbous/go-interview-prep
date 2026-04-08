# Interview Prep Checklist

Ordered by importance for a Tesla Linux Update Systems internship.
Work top-to-bottom. Items higher up have the most impact on your interview.

## Tier 1 — Must Do (core skills they will test)

- [ ] **01 Race Condition** — Go concurrency fundamentals (mutex/atomic, WaitGroup)
- [ ] **02 Deadlock** — Lock ordering, a classic concurrency interview question
- [ ] **05 WaitGroup Bug** — WaitGroup misuse is a very common Go gotcha
- [ ] **04 Channel Panic** — Channel lifecycle (close once, WaitGroup+close pattern)
- [ ] **10 Graceful Shutdown** — Directly relevant: update daemons need clean shutdown
- [ ] **14 C Pthread Race** — C concurrency: pthreads, mutexes, variable capture
- [ ] **Design Q2: Atomic Updates & Rollback (A/B partitions)** — Core to OTA systems
- [ ] **Design Q3: Secure Update Delivery** — Code signing is listed in the JD

## Tier 2 — High Value (likely to come up)

- [ ] **06 Concurrent Map** — Common Go runtime panic, maps + concurrency
- [ ] **07 Context Cancellation** — Context is everywhere in Go systems code
- [ ] **03 Goroutine Leak** — Buffered channels, preventing resource leaks
- [ ] **13 C Use-After-Free** — Linked list bugs, pointer safety in C
- [ ] **11 C Buffer Overflow** — String handling in C, bounds checking
- [ ] **Design Q4: Updating MCUs via Linux Host** — CAN/SPI, the JD mentions peripherals
- [ ] **Design Q7: Concurrent Update Orchestration in Go** — DAG execution, Go concurrency

## Tier 3 — Solid Prep (differentiators)

- [ ] **15 HTTP Download with Resume** — Range requests, checksum verification
- [ ] **17 Process Execution** — os/exec with timeouts, relevant to running install scripts
- [ ] **18 File Locking** — flock(2), preventing concurrent updaters
- [ ] **12 C Memory Leak** — Error path memory management in C
- [ ] **Design Q1: End-to-End OTA Pipeline** — Big picture system design
- [ ] **Design Q5: Fleet Rollout Strategy & Metrics** — Directly in the JD responsibilities

## Tier 4 — If You Have Time (bonus knowledge)

- [ ] **08 Atomic File Write** — Temp+rename pattern for crash safety
- [ ] **09 Slice Gotcha** — Go slice backing array, good to know
- [ ] **16 TCP Server** — Connection handling, resource management
- [ ] **Design Q6: Delta Updates** — Bandwidth optimization
- [ ] **Design Q8: Debugging Fleet-Wide Failure** — Troubleshooting methodology
- [ ] **Design Q9: USB/Ethernet Local Update** — Factory/service workflows
- [ ] **Design Q10: Rate Limiting Downloads** — Scale and CDN design

## How to Use This List

1. **Coding exercises**: Fix the bugs, run `go test -race` or `make test`.
   Aim to complete each in 15-25 minutes (interview pace).
2. **Design questions**: Practice talking through your answer out loud for
   5-10 minutes. Then check the hints. Repeat until fluent.
3. **When done**: Paste your solutions + `GRADING.md` into an AI and ask
   it to grade your work.
4. **Time-box**: If stuck on an exercise for >30 min, read the grading
   rubric for hints, then try again.
