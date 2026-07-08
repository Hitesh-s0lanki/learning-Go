package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	Name     string
	Price    decimal.Decimal
	Quantity int
}

func calculateTotal(items []OrderItem) decimal.Decimal {
	total := decimal.NewFromInt(0)

	for _, item := range items {
		lineTotal := item.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
		total = total.Add(lineTotal)
	}

	return total
}

func main() {
	fmt.Println("----- External Modules Example -----")

	// github.com/google/uuid is an external GitHub module.
	orderID := uuid.NewString()
	fmt.Println("Order ID:", orderID)

	// github.com/shopspring/decimal is useful for money values.
	items := []OrderItem{
		{Name: "T-Shirt", Price: decimal.NewFromFloat(499.99), Quantity: 2},
		{Name: "Mug", Price: decimal.NewFromFloat(199.50), Quantity: 1},
		{Name: "Sticker Pack", Price: decimal.NewFromFloat(99.00), Quantity: 3},
	}

	for _, item := range items {
		fmt.Printf("%s x %d = Rs. %s\n",
			item.Name,
			item.Quantity,
			item.Price.Mul(decimal.NewFromInt(int64(item.Quantity))).StringFixed(2),
		)
	}

	total := calculateTotal(items)
	fmt.Println("Total:", total.StringFixed(2))
}
