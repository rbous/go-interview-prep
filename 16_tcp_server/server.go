package tcp_server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// UpdateStatusServer is a simple TCP server that lets clients query
// update status. Protocol:
//   Client sends: "STATUS <package_name>\n"
//   Server sends: "<package_name>: <status>\n"
//   Client sends: "QUIT\n" to disconnect.
//
// BUG(1): Accepted connections are never closed — resource leak when clients
//         disconnect without sending QUIT.
// BUG(2): The listener is not closed on Shutdown, so Accept blocks forever.
// BUG(3): No connection timeout — a slow/idle client holds a goroutine forever.
// BUG(4): connCount is modified without synchronization.

type UpdateStatusServer struct {
	statuses  map[string]string
	mu        sync.RWMutex
	listener  net.Listener
	connCount int
}

func NewUpdateStatusServer() *UpdateStatusServer {
	return &UpdateStatusServer{
		statuses: make(map[string]string),
	}
}

func (s *UpdateStatusServer) SetStatus(pkg, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statuses[pkg] = status
}

func (s *UpdateStatusServer) GetStatus(pkg string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if st, ok := s.statuses[pkg]; ok {
		return st
	}
	return "unknown"
}

// Start begins listening on the given address. Returns the actual address
// (useful when port is 0).
func (s *UpdateStatusServer) Start(addr string) (string, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", err
	}
	s.listener = ln

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			s.connCount++
			go s.handleConn(conn)
		}
	}()

	return ln.Addr().String(), nil
}

func (s *UpdateStatusServer) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "QUIT" {
			return
		}

		if strings.HasPrefix(line, "STATUS ") {
			pkg := strings.TrimPrefix(line, "STATUS ")
			status := s.GetStatus(pkg)
			fmt.Fprintf(conn, "%s: %s\n", pkg, status)
		} else {
			fmt.Fprintf(conn, "ERR unknown command\n")
		}
	}
}

// Shutdown stops the server.
func (s *UpdateStatusServer) Shutdown() {
	// BUG: doesn't close the listener or existing connections
}

// ConnCount returns the number of connections handled.
func (s *UpdateStatusServer) ConnCount() int {
	return s.connCount
}
