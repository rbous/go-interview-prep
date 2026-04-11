package composition

import "strings"

// Middleware composition exercise.
//
// You're building an HTTP-like middleware pipeline using interface composition
// and the decorator pattern. Each middleware wraps a Handler and adds behavior.
//
// Bugs to fix:
// - AuthMiddleware doesn't properly delegate to the next handler.
// - LogMiddleware has incorrect composition (doesn't embed or wrap correctly).
// - BuildPipeline chains the middlewares in the wrong order.
// - Response field is never populated correctly in one middleware.
//
// Rules:
// - Do NOT modify the test file.

// Request represents a simplified HTTP request.
type Request struct {
	Path    string
	Headers map[string]string
	Body    string
}

// Response represents a simplified HTTP response.
type Response struct {
	Status int
	Body   string
}

// Handler processes a request and returns a response.
type Handler interface {
	Handle(req Request) Response
}

// --- Base handler ---

// EchoHandler returns the request body as the response. This is the "real" handler.
type EchoHandler struct{}

func (h EchoHandler) Handle(req Request) Response {
	return Response{Status: 200, Body: "echo: " + req.Body}
}

// --- Middlewares ---

// AuthMiddleware rejects requests without an "Authorization" header.
type AuthMiddleware struct {
	Next Handler
}

func (a *AuthMiddleware) Handle(req Request) Response {
	if req.Headers == nil || req.Headers["Authorization"] == "" {
		return Response{Status: 401, Body: "unauthorized"}
	}
	// BUG: should delegate to a.Next, but returns an empty response instead
	return Response{Status: 200}
}

// LogMiddleware records each request path it sees, then delegates.
type LogMiddleware struct {
	Next    Handler
	Entries []string
}

func (l *LogMiddleware) Handle(req Request) Response {
	l.Entries = append(l.Entries, req.Path)
	return l.Next.Handle(req)
}

// UppercaseMiddleware transforms the response body to uppercase.
type UppercaseMiddleware struct {
	Next Handler
}

func (u *UppercaseMiddleware) Handle(req Request) Response {
	resp := u.Next.Handle(req)
	resp.Body = strings.ToUpper(resp.Body)
	return resp
}

// BuildPipeline assembles the middleware chain.
// The expected order is: Log -> Auth -> Uppercase -> EchoHandler
// So a request flows: Log first, then Auth check, then through Uppercase, then EchoHandler.
//
// The LogMiddleware instance is returned separately so tests can inspect its log entries.
func BuildPipeline() (Handler, *LogMiddleware) {
	echo := EchoHandler{}
	upper := &UppercaseMiddleware{Next: echo}
	auth := &AuthMiddleware{Next: upper}
	logger := &LogMiddleware{Next: echo} // BUG: should wrap auth, not echo

	return logger, logger
}
