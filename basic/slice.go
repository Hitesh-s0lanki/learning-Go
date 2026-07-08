package main

import "fmt"

func main() {
	// 1. Create a slice.
	numbers := []int{10, 20, 30, 40, 50}
	fmt.Println("Numbers:", numbers)

	// 2. Create a slice from an array.
	arr := [5]string{"Apple", "Banana", "Mango", "Orange", "Grapes"}
	fruits := arr[1:4]
	fmt.Println("Slice from array:", fruits)

	// 3. Create a slice using make.
	scores := make([]int, 3)
	scores[0] = 80
	scores[1] = 90
	scores[2] = 75
	fmt.Println("Scores:", scores)

	// 4. Length and capacity.
	fmt.Println("Length:", len(numbers))
	fmt.Println("Capacity:", cap(numbers))

	// 5. Append single value.
	numbers = append(numbers, 60)
	fmt.Println("After append:", numbers)

	// 6. Append multiple values.
	numbers = append(numbers, 70, 80, 90)
	fmt.Println("After multiple append:", numbers)

	// 7. Slice ranges.
	fmt.Println("First three:", numbers[:3])
	fmt.Println("From index two:", numbers[2:])
	fmt.Println("Middle values:", numbers[2:5])

	// 8. Update slice value.
	numbers[0] = 100
	fmt.Println("After update:", numbers)

	// 9. Loop through slice using range.
	for index, value := range numbers {
		fmt.Println("Index:", index, "Value:", value)
	}

	// 10. Copy a slice.
	copiedNumbers := make([]int, len(numbers))
	copy(copiedNumbers, numbers)
	fmt.Println("Copied numbers:", copiedNumbers)

	// 11. Delete an item from a slice.
	// Delete value at index 2.
	indexToDelete := 2
	numbers = append(numbers[:indexToDelete], numbers[indexToDelete+1:]...)
	fmt.Println("After delete:", numbers)

	// 12. Search in a slice.
	searchValue := 70
	found := false
	for _, value := range numbers {
		if value == searchValue {
			found = true
			break
		}
	}
	fmt.Println("Is 70 available?", found)

	// 13. Empty slice.
	var emptySlice []string
	fmt.Println("Empty slice:", emptySlice)
	fmt.Println("Is empty slice nil?", emptySlice == nil)
}
