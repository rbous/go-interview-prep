package deadlock

import "sync"

// Account represents a bank account with a balance protected by a mutex.
// Transfer should move `amount` from account `a` to account `b` atomically.
//
// BUG: When two goroutines call Transfer(a, b, 10) and Transfer(b, a, 5)
// simultaneously, a deadlock can occur. Fix the locking strategy so
// transfers never deadlock, while remaining safe for concurrent use.

type Account struct {
	mu      sync.Mutex
	id      int
	balance int
}

func NewAccount(id, balance int) *Account {
	return &Account{id: id, balance: balance}
}

func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func Transfer(from, to *Account, amount int) {
	from.mu.Lock()
	defer from.mu.Unlock()

	to.mu.Lock()
	defer to.mu.Unlock()

	if from.balance >= amount {
		from.balance -= amount
		to.balance += amount
	}
}
