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

type Transaction struct {
	FromId string
	ToId   string
	Amount float64
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (p *PaymentSystem) AddUser(user *User) {
	p.Users[user.ID] = user
}

func (p *PaymentSystem) AddTransaction(transaction Transaction) {
	p.TransactionQueue = append(p.TransactionQueue, transaction)
}

func (u *User) Deposit(value float64) {
	fmt.Println("Пополнение пользователю", u.Name, "на", value)
	u.Balance += value
}

func (u *User) Withdraw(value float64) error {
	if u.Balance-value < 0 {
		return errors.New("Недостаточно средств на балансе")
	}

	fmt.Println("Снятие со счета пользователя", u.Name, value)
	u.Balance -= value
	return nil
}

func (p *PaymentSystem) ProcessingTransactions(transaction Transaction) error {
	userTo, ok := p.Users[transaction.ToId]
	if ok == false {
		return errors.New("Пользоватьель ToId не найден")
	}

	userFrom, ok := p.Users[transaction.FromId]
	if ok == false {
		return errors.New("Пользоватьель FromId не найден")
	}

	fmt.Println("Перевожу с пользователя ID:", transaction.FromId, "на  ID:",
		transaction.ToId, "сумму в размере", transaction.Amount)

	err := userFrom.Withdraw(transaction.Amount)
	if err != nil {
		return err
	}

	userTo.Deposit(transaction.Amount)

	return nil
}

func main() {
	ps := &PaymentSystem{Users: make(map[string]*User), TransactionQueue: make([]Transaction, 0)}

	user1 := &User{ID: "id1", Name: "user1", Balance: 100.0}
	user2 := &User{ID: "id2", Name: "user2", Balance: 500.0}

	ps.AddUser(user1)
	ps.AddUser(user2)

	tr1 := Transaction{FromId: "id1", ToId: "id2", Amount: 200}
	tr2 := Transaction{FromId: "id2", ToId: "id1", Amount: 50}

	ps.AddTransaction(tr1)
	ps.AddTransaction(tr2)

	for i := 0; i < len(ps.TransactionQueue); i++ {
		err := ps.ProcessingTransactions(ps.TransactionQueue[i])
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("Баланс", user1.Name, user1.Balance)
	fmt.Println("Баланс", user2.Name, user2.Balance)
}
