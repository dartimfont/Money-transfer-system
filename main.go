package main

import (
	"errors"
	"fmt"
	"sync"
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

func (ps *PaymentSystem) AddUser(user *User) {
	ps.Users[user.ID] = user
}

func (ps *PaymentSystem) AddTransaction(transaction Transaction) {
	ps.TransactionQueue = append(ps.TransactionQueue, transaction)
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

func (ps *PaymentSystem) ProcessingTransactions(transaction Transaction) error {
	userTo, ok := ps.Users[transaction.ToId]
	if ok == false {
		return errors.New("Пользоватьель ToId не найден")
	}

	userFrom, ok := ps.Users[transaction.FromId]
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

func (ps *PaymentSystem) worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	for tr := range ch {
		err := ps.ProcessingTransactions(tr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	ps := &PaymentSystem{Users: make(map[string]*User), TransactionQueue: make([]Transaction, 0)}

	user1 := &User{ID: "id1", Name: "user1", Balance: 1000.0}
	user2 := &User{ID: "id2", Name: "user2", Balance: 500.0}

	ps.AddUser(user1)
	ps.AddUser(user2)

	tr1 := Transaction{FromId: "id1", ToId: "id2", Amount: 200}
	tr2 := Transaction{FromId: "id2", ToId: "id1", Amount: 50}

	ps.AddTransaction(tr1)
	ps.AddTransaction(tr2)

	ch := make(chan Transaction, 3)
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go ps.worker(ch, &wg)
	}

	for _, tr := range ps.TransactionQueue {
		ch <- tr
	}

	close(ch)
	wg.Wait()

	fmt.Println("Баланс", user1.Name, user1.Balance)
	fmt.Println("Баланс", user2.Name, user2.Balance)
}
