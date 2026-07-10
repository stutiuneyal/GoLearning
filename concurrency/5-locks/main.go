package main

import (
	"fmt"
	"sync"
	"time"
)

type BankAccount struct {
	balance int
	mu      sync.Mutex
}

// Deposit
func (b *BankAccount) Deposit(amount int) {

	b.mu.Lock()
	defer b.mu.Unlock()

	//CS
	b.balance += amount

	fmt.Println("Deposited = ", amount)

}

// Withdraw
func (b *BankAccount) Withdraw(amount int) {

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.balance < amount {
		fmt.Println("Cannot withdraw = ", amount)
		return
	}

	// CS
	b.balance -= amount

	fmt.Println("Withdraw = ", amount)
}

// Balance
func (b *BankAccount) Balance() int {

	b.mu.Lock()
	defer b.mu.Unlock()

	return b.balance
}

func main() {

	var wg sync.WaitGroup

	account := &BankAccount{
		balance: 100,
	}

	// launch deposit goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(amount int) {
			defer wg.Done()

			// simulate work
			time.Sleep(time.Duration(amount) * time.Millisecond)

			account.Deposit(amount)
			fmt.Println("Current Balance: ", account.Balance())

		}(i + 1)
	}

	// launch withdraw go routines
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(amount int) {
			defer wg.Done()

			// simulate work
			time.Sleep(time.Duration(amount) * time.Millisecond)

			account.Withdraw(amount * 10)
			fmt.Println("Current Balance: ", account.Balance())

		}(i + 1)
	}

	wg.Wait()

	fmt.Println(account.Balance())
}
