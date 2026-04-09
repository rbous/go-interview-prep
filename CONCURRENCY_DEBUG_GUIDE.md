# Go Concurrency: Interview Debugging Guide

When handed a concurrent Go program in an interview, do not instantly read the code line-by-line. Follow this structured approach to demonstrate seniority, logic, and a deep understanding of Go's concurrency model.

## 1. The 3-Step Mental Framework

> [!TIP]
> Think about goroutines on a shared timeline, not in isolation. A delay in one goroutine directly blocks another.

1. **Map the Architecture**: Identify the "Actors" (Main thread, Workers, Producers) and the "Synchronization Primitives" (Channels, WaitGroups, Mutexes).
2. **State Your Suspicions**: Before deep-diving, explicitly state what usually breaks in this architecture (e.g., *"Since we have unbuffered channels, I'm going to look out for send/receive mismatched timing"*).
3. **Trace the Timeline**: Read top-to-bottom synchronously. When you hit a blocking operation (like `wg.Wait()` or `ch <- val`), freeze the timeline and ask: *"At this exact millisecond, what is every other goroutine doing? Are they capable of unblocking me?"*

## 2. The "Gotcha" Checklist

Scan the code for these 6 common anti-patterns:

- [ ] **The Producer/Consumer Deadlock**: An unbuffered channel operation (`ch <- val` or `<-ch`) without a concurrent goroutine ready to interact with it. Look for `wg.Wait()` barriers placed between senders and receivers.
- [ ] **Goroutine Leaks (Memory Leaks)**: A goroutine blocked on a channel receive (`range ch`) because the producer never called `close(ch)`.
- [ ] **Panic on Closed Channel**: Multiple producers trying to close the same channel, or writing to a channel that a coordinator already closed. *(Rule: Only the sender closes. If many senders, a coordinator handles it).*
- [ ] **Pass-by-Value WaitGroups**: Passing `sync.WaitGroup` to a function without a pointer (`*sync.WaitGroup`), causing it to copy the lock and break `Wait()`.
- [ ] **The Loop Variable Trap**: Launching goroutines inside a `for i := 0` loop and using `i` directly, causing all goroutines to share the final loop value.
- [ ] **Concurrent Mutations (Data Races)**: Multiple goroutines writing to a map or slice without a `sync.Mutex`. Go maps are explicitly not thread-safe.

## 3. Practical Tooling

> [!IMPORTANT]
> The Go test runner provides the most powerful tool for debugging deadlocks built directly into the language.

### Reading the Panic Trace
Run `go test -timeout 5s`. When the timeout hits, Go dumps the state of every alive goroutine.
1. **Count the living:** See which actors never exited.
2. **Read the state tag:**
   - `[sync.WaitGroup.Wait]`: Stuck waiting for `Done()` calls.
   - `[chan receive]`: Stuck waiting for someone to send on a channel.
   - `[chan send]`: Stuck waiting for someone to read from a channel.
3. **Find the Knot:** Match the goroutine stuck on `Wait` against the goroutine stuck on a `chan` operation to pinpoint the deadlock.

### The "Low-Tech" Trace
If the panic trace isn't clear, scatter `fmt.Println` statements around blocking operations with clear prefix labels:
```go
fmt.Println("[Worker] Sending Result")
resultCh <- val
fmt.Println("[Worker] Result Sent Successfully") // If this never prints, you found the block.
```
