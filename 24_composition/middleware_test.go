package composition

import "testing"

func authedRequest(path, body string) Request {
	return Request{
		Path:    path,
		Headers: map[string]string{"Authorization": "Bearer token123"},
		Body:    body,
	}
}

func TestEchoHandler(t *testing.T) {
	h := EchoHandler{}
	resp := h.Handle(Request{Body: "hello"})
	if resp.Body != "echo: hello" {
		t.Errorf("EchoHandler: got %q, want %q", resp.Body, "echo: hello")
	}
}

func TestAuthMiddlewareRejectsUnauthorized(t *testing.T) {
	pipeline, _ := BuildPipeline()
	resp := pipeline.Handle(Request{Path: "/secret", Body: "data"})

	if resp.Status != 401 {
		t.Errorf("expected 401 for unauthorized request, got %d", resp.Status)
	}
	if resp.Body != "unauthorized" {
		t.Errorf("expected 'unauthorized' body, got %q", resp.Body)
	}
}

func TestAuthMiddlewareDelegates(t *testing.T) {
	pipeline, _ := BuildPipeline()
	resp := pipeline.Handle(authedRequest("/api", "world"))

	if resp.Status != 200 {
		t.Errorf("expected 200 for authorized request, got %d", resp.Status)
	}
	if resp.Body != "ECHO: WORLD" {
		t.Errorf("expected uppercased echo, got %q", resp.Body)
	}
}

func TestLogMiddlewareRecords(t *testing.T) {
	pipeline, logger := BuildPipeline()

	pipeline.Handle(authedRequest("/first", "a"))
	pipeline.Handle(authedRequest("/second", "b"))

	if len(logger.Entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(logger.Entries))
	}
	if logger.Entries[0] != "/first" || logger.Entries[1] != "/second" {
		t.Errorf("unexpected log entries: %v", logger.Entries)
	}
}

func TestFullPipelineOrder(t *testing.T) {
	pipeline, logger := BuildPipeline()

	// Authorized request should flow: Log -> Auth -> Uppercase -> Echo
	resp := pipeline.Handle(authedRequest("/test", "go"))

	// Logger should have recorded the path
	if len(logger.Entries) != 1 || logger.Entries[0] != "/test" {
		t.Errorf("logger didn't record path: %v", logger.Entries)
	}

	// Response should be uppercased echo
	if resp.Body != "ECHO: GO" {
		t.Errorf("pipeline output = %q, want %q", resp.Body, "ECHO: GO")
	}
}

func TestUnauthorizedStillLogged(t *testing.T) {
	pipeline, logger := BuildPipeline()

	resp := pipeline.Handle(Request{Path: "/blocked", Body: "x"})

	// Even unauthorized requests should be logged
	if len(logger.Entries) != 1 || logger.Entries[0] != "/blocked" {
		t.Errorf("unauthorized request not logged: %v", logger.Entries)
	}
	if resp.Status != 401 {
		t.Errorf("expected 401, got %d", resp.Status)
	}
}
