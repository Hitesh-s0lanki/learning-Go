package main

import (
	"fmt"

	"go-learning/custommodule/calculator"
)

func main() {
	fmt.Println("----- Custom Module Example -----")

	sum := calculator.Add(10, 20)
	difference := calculator.Subtract(50, 15)

	fmt.Println("Sum:", sum)
	fmt.Println("Difference:", difference)

	result, err := calculator.Divide(100, 5)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Division result:", result)
	}

	items := []calculator.InvoiceItem{
		{Name: "T-Shirt", Price: 499.99, Quantity: 2},
		{Name: "Mug", Price: 199.50, Quantity: 1},
		{Name: "Sticker Pack", Price: 99.00, Quantity: 3},
	}

	total := calculator.CalculateTotal(items)
	fmt.Printf("Invoice total: %.2f\n", total)
}
