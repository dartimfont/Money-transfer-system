package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID      string
	Name    string
	Balance float64
}

func (u *User) Deposit(value float64) {
	fmt.Println("Пополнение пользователю", u.Name, "на", value)
	u.Balance += value
	u.printBalance()
}

func (u *User) Withdraw(value float64) error {
	fmt.Println("Снятие со счета пользователя", u.Name, value)

	if u.Balance-value < 0 {
		return errors.New("Недостаточно средств на балансе")
	}

	u.Balance -= value
	u.printBalance()
	return nil
}

func (u *User) printBalance() {
	fmt.Println("Баланс", u.Name, u.Balance)
}

func main() {
	user1 := &User{ID: "id1", Name: "user1", Balance: 12.0}
	user2 := &User{ID: "id2", Name: "user2", Balance: 0.0}
	user3 := &User{ID: "id3", Name: "user3", Balance: 45.0}

	user1.Deposit(30)
	user2.Deposit(30)
	user3.Deposit(30)

	if err := user1.Withdraw(100); err != nil {
		fmt.Println(err.Error())
	}

	if err := user2.Withdraw(10); err != nil {
		fmt.Println(err.Error())
	}

	if err := user3.Withdraw(560); err != nil {
		fmt.Println(err.Error())
	}
}
