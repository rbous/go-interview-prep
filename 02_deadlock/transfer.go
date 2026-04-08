package deadlock

import "sync"

// Account represents a bank account with a balance protected by a mutex.
// Transfer should move `amount` from account `a` to account `b` atomically.

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

	if from.balance >= amount {
		from.balance -= amount
		from.mu.Unlock()
		to.mu.Lock()
		to.balance += amount
		to.mu.Unlock()
	} else {
		from.mu.Unlock()
	}
}
