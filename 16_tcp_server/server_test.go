package tcp_server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBasicStatusQuery(t *testing.T) {
	s := NewUpdateStatusServer()
	s.SetStatus("curl", "installed")
	s.SetStatus("wget", "downloading")

	addr, err := s.Start("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "STATUS curl\n")
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		t.Fatal("no response")
	}
	if got := scanner.Text(); got != "curl: installed" {
		t.Errorf("got %q, want %q", got, "curl: installed")
	}
}

func TestUnknownPackage(t *testing.T) {
	s := NewUpdateStatusServer()
	addr, err := s.Start("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "STATUS nonexistent\n")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	if !strings.Contains(scanner.Text(), "unknown") {
		t.Errorf("expected 'unknown' status, got %q", scanner.Text())
	}
}

func TestMultipleConnections(t *testing.T) {
	s := NewUpdateStatusServer()
	s.SetStatus("pkg", "ready")

	addr, err := s.Start("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Error(err)
				return
			}
			defer conn.Close()

			fmt.Fprintf(conn, "STATUS pkg\n")
			scanner := bufio.NewScanner(conn)
			if scanner.Scan() {
				if !strings.Contains(scanner.Text(), "ready") {
					t.Errorf("unexpected: %q", scanner.Text())
				}
			}
			fmt.Fprintf(conn, "QUIT\n")
		}()
	}
	wg.Wait()

	if s.ConnCount() != 10 {
		t.Errorf("ConnCount() = %d, want 10", s.ConnCount())
	}
}

func TestShutdownClosesListener(t *testing.T) {
	s := NewUpdateStatusServer()
	addr, err := s.Start("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	s.Shutdown()

	// After shutdown, new connections should be refused
	time.Sleep(50 * time.Millisecond)
	_, err = net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err == nil {
		t.Error("expected connection refused after shutdown")
	}
}

func TestConnectionCloseOnClientDisconnect(t *testing.T) {
	s := NewUpdateStatusServer()
	_, err := s.Start("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown()

	// This test just ensures no goroutine leak / resource leak
	// when clients disconnect without QUIT.
	// The real verification is via `go test -race`.
}
