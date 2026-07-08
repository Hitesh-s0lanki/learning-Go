package main

import "fmt"

func main() {
	// Map example: student name -> marks
	studentMarks := map[string]int{
		"Hemant": 85,
		"Amit":   72,
		"Neha":   91,
	}

	studentName := "Hemant"

	// Conditional statement using map lookup.
	// ok is true if the key exists in the map.
	if marks, ok := studentMarks[studentName]; ok {
		fmt.Println(studentName, "marks:", marks)

		if marks >= 90 {
			fmt.Println("Grade: A")
		} else if marks >= 75 {
			fmt.Println("Grade: B")
		} else if marks >= 50 {
			fmt.Println("Grade: C")
		} else {
			fmt.Println("Grade: Fail")
		}
	} else {
		fmt.Println(studentName, "not found in the map")
	}

	// Another condition: checking a missing key.
	searchName := "Rahul"
	if marks, ok := studentMarks[searchName]; ok {
		fmt.Println(searchName, "marks:", marks)
	} else {
		fmt.Println(searchName, "not found in the map")
	}
}
