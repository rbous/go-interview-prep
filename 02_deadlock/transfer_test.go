package deadlock

import (
	"sync"
	"testing"
)

func TestTransferBasic(t *testing.T) {
	a := NewAccount(1, 100)
	b := NewAccount(2, 50)

	Transfer(a, b, 30)

	if a.Balance() != 70 {
		t.Errorf("a.Balance() = %d, want 70", a.Balance())
	}
	if b.Balance() != 80 {
		t.Errorf("b.Balance() = %d, want 80", b.Balance())
	}
}

func TestTransferInsufficientFunds(t *testing.T) {
	a := NewAccount(1, 10)
	b := NewAccount(2, 50)

	Transfer(a, b, 20)

	if a.Balance() != 10 {
		t.Errorf("a.Balance() = %d, want 10", a.Balance())
	}
	if b.Balance() != 50 {
		t.Errorf("b.Balance() = %d, want 50", b.Balance())
	}
}

func TestTransferConcurrentBidirectional(t *testing.T) {
	a := NewAccount(1, 1000)
	b := NewAccount(2, 1000)

	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			Transfer(a, b, 1)
		}()
		go func() {
			defer wg.Done()
			Transfer(b, a, 1)
		}()
	}
	wg.Wait()

	totalBefore := 2000
	totalAfter := a.Balance() + b.Balance()
	if totalAfter != totalBefore {
		t.Errorf("total balance = %d, want %d (money created or destroyed)", totalAfter, totalBefore)
	}
}
