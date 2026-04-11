package interface_composition

import "testing"

func TestPrefixNotifierSendsFormattedMessage(t *testing.T) {
    n := NewPrefixNotifier("INFO")

    if err := n.Notify("system started"); err != nil {
        t.Fatalf("Notify returned unexpected error: %v", err)
    }

    sent := n.Sent()
    if len(sent) != 1 {
        t.Fatalf("got %d sent messages, want 1", len(sent))
    }
    if want := "INFO: system started"; sent[0] != want {
        t.Errorf("sent[0] = %q, want %q", sent[0], want)
    }
}

func TestPrefixNotifierMultipleMessages(t *testing.T) {
    n := NewPrefixNotifier("WARN")

    n.Notify("low memory")
    n.Notify("disk usage high")

    sent := n.Sent()
    if len(sent) != 2 {
        t.Fatalf("got %d sent messages, want 2", len(sent))
    }

    want := []string{"WARN: low memory", "WARN: disk usage high"}
    for i, w := range want {
        if sent[i] != w {
            t.Errorf("sent[%d] = %q, want %q", i, sent[i], w)
        }
    }
}

func TestPrefixNotifierIsolation(t *testing.T) {
    // Two notifiers should not share state.
    info := NewPrefixNotifier("INFO")
    warn := NewPrefixNotifier("WARN")

    info.Notify("started")
    warn.Notify("low battery")

    if got := info.Sent()[0]; got != "INFO: started" {
        t.Errorf("info: got %q, want %q", got, "INFO: started")
    }
    if got := warn.Sent()[0]; got != "WARN: low battery" {
        t.Errorf("warn: got %q, want %q", got, "WARN: low battery")
    }
}
