package main

import "fmt"

func main() {
	day := "Monday"

	// Basic switch case example.
	switch day {
	case "Monday":
		fmt.Println("Start of the week")
	case "Tuesday", "Wednesday", "Thursday":
		fmt.Println("Middle of the week")
	case "Friday":
		fmt.Println("Almost weekend")
	case "Saturday", "Sunday":
		fmt.Println("Weekend")
	default:
		fmt.Println("Invalid day")
	}

	marks := 85

	// Switch without an expression.
	// This is useful when checking conditions.
	switch {
	case marks >= 90:
		fmt.Println("Grade: A")
	case marks >= 75:
		fmt.Println("Grade: B")
	case marks >= 50:
		fmt.Println("Grade: C")
	default:
		fmt.Println("Grade: Fail")
	}
}
