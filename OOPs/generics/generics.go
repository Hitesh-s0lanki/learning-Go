package main

import "fmt"

type Number interface {
	int | int64 | float64
}

type Pair[T any] struct {
	First  T
	Second T
}

func PrintItems[T any](items []T) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func Add[T Number](a T, b T) T {
	return a + b
}

func FindIndex[T comparable](items []T, searchValue T) int {
	for index, item := range items {
		if item == searchValue {
			return index
		}
	}

	return -1
}

func main() {
	fmt.Println("----- Generic Function With Any Type -----")

	names := []string{"Hemant", "Amit", "Neha"}
	PrintItems(names)

	numbers := []int{10, 20, 30}
	PrintItems(numbers)

	fmt.Println("----- Generic Function With Number Constraint -----")

	fmt.Println("Int add:", Add(10, 20))
	fmt.Println("Float add:", Add(10.5, 20.25))

	fmt.Println("----- Generic Function With Comparable Constraint -----")

	fmt.Println("Index of Amit:", FindIndex(names, "Amit"))
	fmt.Println("Index of 30:", FindIndex(numbers, 30))
	fmt.Println("Index of 100:", FindIndex(numbers, 100))

	fmt.Println("----- Generic Struct -----")

	intPair := Pair[int]{First: 10, Second: 20}
	stringPair := Pair[string]{First: "Go", Second: "Generics"}

	fmt.Println("Int pair:", intPair)
	fmt.Println("String pair:", stringPair)
}
