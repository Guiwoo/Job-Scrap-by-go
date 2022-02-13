package banking

import (
	"errors"
	"fmt"
)

// BankAcc Strut
type BankAcc struct {
	owner   string
	balance int
}

var errNoMoney = errors.New("can not withdraw")

func NewAcc(owner string) *BankAcc {
	acc := BankAcc{owner: owner, balance: 0}
	return &acc
}

//method reciver
func (b *BankAcc) Deposit(amount int) {
	b.balance += amount
}

func (b BankAcc) Balance() int {
	return b.balance
}

func (b *BankAcc) Withdraw(amout int) error {
	if b.balance < amout {
		return errNoMoney
	}
	b.balance -= amout
	return nil
}

func (b *BankAcc) Change(newOwner string) {
	b.owner = newOwner
}

func (b BankAcc) Owner() string {
	return b.owner
}

func (b BankAcc) String() string {
	return fmt.Sprint(b.Owner(), "s account\nHas:", b.Balance())
}
