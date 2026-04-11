package abstraction

import "testing"

func TestMemoryStorePutGet(t *testing.T) {
	s := NewMemoryStore()
	s.Put("host", "localhost")

	val, ok := s.Get("host")
	if !ok {
		t.Fatal("expected key 'host' to exist")
	}
	if val != "localhost" {
		t.Errorf("got %q, want %q", val, "localhost")
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	s := NewMemoryStore()
	s.Put("key", "value")
	s.Delete("key")

	_, ok := s.Get("key")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestLoggedStoreWrapsAnyStore(t *testing.T) {
	// LoggedStore must accept any Store, not just MemoryStore.
	var base Store = NewMemoryStore()
	logged := NewLoggedStore(base)

	logged.Put("env", "production")
	val, ok := logged.Get("env")
	if !ok || val != "production" {
		t.Errorf("Get through LoggedStore failed: got %q, ok=%v", val, ok)
	}
}

func TestLoggedStoreDeleteDelegates(t *testing.T) {
	base := NewMemoryStore()
	logged := NewLoggedStore(base)

	logged.Put("temp", "data")
	logged.Delete("temp")

	_, ok := logged.Get("temp")
	if ok {
		t.Error("Delete through LoggedStore did not delegate to inner store")
	}
}

func TestLoggedStoreRecordsLog(t *testing.T) {
	base := NewMemoryStore()
	logged := NewLoggedStore(base)

	logged.Put("a", "1")
	logged.Get("a")
	logged.Delete("a")

	log := logged.Log()
	if len(log) != 3 {
		t.Fatalf("expected 3 log entries, got %d: %v", len(log), log)
	}

	want := []string{"PUT a=1", "GET a", "DELETE a"}
	for i, w := range want {
		if log[i] != w {
			t.Errorf("log[%d] = %q, want %q", i, log[i], w)
		}
	}
}

func TestGetOrDefault(t *testing.T) {
	s := NewMemoryStore()
	s.Put("exists", "yes")

	if got := GetOrDefault(s, "exists", "no"); got != "yes" {
		t.Errorf("existing key: got %q, want %q", got, "yes")
	}

	if got := GetOrDefault(s, "missing", "fallback"); got != "fallback" {
		t.Errorf("missing key: got %q, want %q", got, "fallback")
	}
}

func TestStoreInterfaceSatisfaction(t *testing.T) {
	// Both MemoryStore and LoggedStore must satisfy the Store interface.
	var _ Store = NewMemoryStore()
	var _ Store = NewLoggedStore(NewMemoryStore())
}
