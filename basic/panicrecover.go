package main

import "fmt"

func divide(a int, b int) int {
	if b == 0 {
		panic("cannot divide by zero")
	}

	return a / b
}

func safeDivide(a int, b int) {
	defer func() {
		errorMessage := recover()
		if errorMessage != nil {
			fmt.Println("Recovered from panic:", errorMessage)
		}
	}()

	result := divide(a, b)
	fmt.Println("Result:", result)
}

func main() {
	fmt.Println("Program started")

	safeDivide(10, 2)
	safeDivide(10, 0)

	fmt.Println("Program continued after recovery")
}
