package calculator

import "errors"

type InvoiceItem struct {
	Name     string
	Price    float64
	Quantity int
}

func Add(a int, b int) int {
	return a + b
}

func Subtract(a int, b int) int {
	return a - b
}

func Divide(a float64, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}

	return a / b, nil
}

func CalculateTotal(items []InvoiceItem) float64 {
	total := 0.0

	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	return total
}
