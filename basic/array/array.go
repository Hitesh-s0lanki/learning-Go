package main

import "fmt"

type Student struct {
	Name string
	Age  int
}

func main() {
	// 1. Declare an array with size.
	// Empty values become zero values: int = 0.
	var numbers [5]int
	fmt.Println("Empty int array:", numbers)

	// 2. Declare and initialize with values.
	var marks [3]int = [3]int{80, 90, 75}
	fmt.Println("Marks:", marks)

	// 3. Short declaration.
	fruits := [3]string{"Apple", "Banana", "Mango"}
	fmt.Println("Fruits:", fruits)

	// 4. Let Go count the array size using ...
	colors := [...]string{"Red", "Green", "Blue", "Yellow"}
	fmt.Println("Colors:", colors)

	// 5. Initialize only selected indexes.
	// Index 0 = 10, index 3 = 40, all other values are 0.
	selectedNumbers := [5]int{0: 10, 3: 40}
	fmt.Println("Selected indexes:", selectedNumbers)

	// 6. Array of booleans.
	flags := [3]bool{true, false, true}
	fmt.Println("Flags:", flags)

	// 7. Array of floats.
	prices := [4]float64{99.99, 149.50, 20.75, 10.00}
	fmt.Println("Prices:", prices)

	// 8. Two-dimensional array.
	matrix := [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println("Matrix:", matrix)

	// 9. Array of structs.
	students := [2]Student{
		{Name: "Hemant", Age: 22},
		{Name: "Amit", Age: 24},
	}
	fmt.Println("Students:", students)

	// 10. Access and update array values.
	fruits[1] = "Orange"
	fmt.Println("Updated fruits:", fruits)
	fmt.Println("First fruit:", fruits[0])
	fmt.Println("Number of fruits:", len(fruits))

	fmt.Println("----- Array Operations -----")

	// 11. Loop through an array using index.
	for i := 0; i < len(marks); i++ {
		fmt.Println("Mark at index", i, ":", marks[i])
	}

	// 12. Loop through an array using range.
	for index, fruit := range fruits {
		fmt.Println("Fruit at index", index, ":", fruit)
	}

	// 13. Find the sum of array values.
	scores := [5]int{10, 20, 30, 40, 50}
	sum := 0
	for _, score := range scores {
		sum += score
	}
	fmt.Println("Sum of scores:", sum)

	// 14. Find the largest value in an array.
	largest := scores[0]
	for _, score := range scores {
		if score > largest {
			largest = score
		}
	}
	fmt.Println("Largest score:", largest)

	// 15. Search for a value in an array.
	searchFruit := "Mango"
	found := false
	for _, fruit := range fruits {
		if fruit == searchFruit {
			found = true
			break
		}
	}
	fmt.Println("Is Mango available?", found)

	// 16. Copy an array.
	// In Go, assigning one array to another creates a copy.
	copiedScores := scores
	copiedScores[0] = 100
	fmt.Println("Original scores:", scores)
	fmt.Println("Copied scores:", copiedScores)

	// 17. Compare arrays.
	// Arrays can be compared if their elements are comparable.
	firstArray := [3]int{1, 2, 3}
	secondArray := [3]int{1, 2, 3}
	fmt.Println("Are arrays equal?", firstArray == secondArray)
}
