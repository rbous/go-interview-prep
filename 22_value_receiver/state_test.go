package value_receiver

import (
    "reflect"
    "testing"
)

func TestApplyChangesVersion(t *testing.T) {
    s := NewUpdateState("v1.0")
    s.Apply("v1.1")

    if got := s.Current(); got != "v1.1" {
        t.Errorf("Current() = %q after Apply, want %q", got, "v1.1")
    }
}

func TestApplyRecordsHistory(t *testing.T) {
    s := NewUpdateState("v1.0")
    s.Apply("v1.1")
    s.Apply("v1.2")

    want := []string{"v1.0", "v1.1"}
    if got := s.History(); !reflect.DeepEqual(got, want) {
        t.Errorf("History() = %v, want %v", got, want)
    }
    if got := s.Current(); got != "v1.2" {
        t.Errorf("Current() = %q, want %q", got, "v1.2")
    }
}

func TestRollbackRestoresPrevious(t *testing.T) {
    s := NewUpdateState("v1.0")
    s.Apply("v1.1")
    s.Apply("v1.2")
    s.Rollback()

    if got := s.Current(); got != "v1.1" {
        t.Errorf("Current() = %q after Rollback, want %q", got, "v1.1")
    }
    want := []string{"v1.0"}
    if got := s.History(); !reflect.DeepEqual(got, want) {
        t.Errorf("History() = %v after Rollback, want %v", got, want)
    }
}

func TestRollbackOnEmptyHistoryIsNoop(t *testing.T) {
    s := NewUpdateState("v1.0")
    s.Rollback() // must not panic

    if got := s.Current(); got != "v1.0" {
        t.Errorf("Current() = %q after empty Rollback, want %q", got, "v1.0")
    }
}
