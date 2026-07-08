package main

import "fmt"

type PaymentMethod interface {
	Pay(amount float64)
	Refund(amount float64)
}

type CreditCard struct {
	cardNumber string
}

type UPI struct {
	upiID string
}

func (card CreditCard) Pay(amount float64) {
	fmt.Printf("Paid %.2f using credit card %s\n", amount, card.cardNumber)
}

func (card CreditCard) Refund(amount float64) {
	fmt.Printf("Refunded %.2f to credit card %s\n", amount, card.cardNumber)
}

func (upi UPI) Pay(amount float64) {
	fmt.Printf("Paid %.2f using UPI ID %s\n", amount, upi.upiID)
}

func (upi UPI) Refund(amount float64) {
	fmt.Printf("Refunded %.2f to UPI ID %s\n", amount, upi.upiID)
}

func checkout(payment PaymentMethod, amount float64) {
	payment.Pay(amount)
}

func returnPayment(payment PaymentMethod, amount float64) {
	payment.Refund(amount)
}

func main() {
	card := CreditCard{cardNumber: "XXXX-XXXX-1234"}
	upi := UPI{upiID: "hemant@upi"}

	checkout(card, 1500.00)
	checkout(upi, 499.99)

	returnPayment(card, 500.00)
	returnPayment(upi, 100.00)
}
