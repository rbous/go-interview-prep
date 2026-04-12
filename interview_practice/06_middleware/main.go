package main

import "fmt"

// Request Pipeline
//
// Middleware should run in registration order: first registered executes first.
// The chain should log, then authenticate, then reach the handler.
//
// Expected output:
//   [log] request: hello
//   [auth] request: hello
//   [core] request: hello

type Handler func(string)

func Chain(h Handler, mws ...func(Handler) Handler) Handler {
    for _, mw := range mws {
        h = mw(h)
    }
    return h
}

func WithLog(next Handler) Handler {
    return func(s string) {
        fmt.Println("[log] request:", s)
        next(s)
    }
}

func WithAuth(next Handler) Handler {
    return func(s string) {
        fmt.Println("[auth] request:", s)
        next(s)
    }
}

func main() {
    core := func(s string) {
        fmt.Println("[core] request:", s)
    }

    h := Chain(core, WithLog, WithAuth)
    h("hello")
}
