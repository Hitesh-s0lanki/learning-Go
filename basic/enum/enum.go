package main

import "fmt"

// Go does not have an enum keyword.
// We create enum-like values using a custom type and iota.
type OrderStatus int

const (
	Pending OrderStatus = iota
	Processing
	Shipped
	Delivered
	Cancelled
)

func main() {
	var status OrderStatus = Shipped

	fmt.Println("Order status value:", status)

	if status == Shipped {
		fmt.Println("Your order has been shipped.")
	}
}
