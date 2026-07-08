package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func updateValueWithoutPointer(value int) {
	value = 100
}

func updateValueWithPointer(value *int) {
	*value = 100
}

func updateUserAge(user *User) {
	user.Age = 25
}

func main() {
	// Pointer means: a variable that stores the memory address of another variable.
	// & gives the address.
	// * gives the value stored at that address.

	age := 22

	fmt.Println("Age value:", age)
	fmt.Println("Age address:", &age)

	// Create a pointer variable.
	var agePointer *int = &age

	fmt.Println("Pointer stores address:", agePointer)
	fmt.Println("Value from pointer:", *agePointer)

	// Update value using pointer.
	*agePointer = 23
	fmt.Println("Updated age:", age)

	fmt.Println("----- Pointer With Function -----")

	number := 50
	updateValueWithoutPointer(number)
	fmt.Println("Without pointer:", number)

	updateValueWithPointer(&number)
	fmt.Println("With pointer:", number)

	fmt.Println("----- Pointer With Struct -----")

	user := User{Name: "Hemant", Age: 22}
	fmt.Println("Before update:", user)

	updateUserAge(&user)
	fmt.Println("After update:", user)

	fmt.Println("----- Pointer With Array -----")

	numbers := [3]int{10, 20, 30}
	numbersPointer := &numbers

	numbersPointer[0] = 100
	fmt.Println("Updated array:", numbers)

	fmt.Println("----- Slice Behavior -----")

	// Slices already point to an underlying array internally.
	// That is why changing a slice inside another variable can affect the original data.
	scores := []int{80, 90, 75}
	otherScores := scores
	otherScores[0] = 100

	fmt.Println("Original scores:", scores)
	fmt.Println("Other scores:", otherScores)

	fmt.Println("----- Nil Pointer -----")

	var emptyPointer *int
	fmt.Println("Empty pointer:", emptyPointer)

	if emptyPointer == nil {
		fmt.Println("Pointer is nil, so do not dereference it")
	}
}
