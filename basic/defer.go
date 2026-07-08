package main

import "fmt"

func openFile() {
	fmt.Println("Opening file")
}

func closeFile() {
	fmt.Println("Closing file")
}

func processFile() {
	openFile()
	defer closeFile()

	fmt.Println("Reading file")
	fmt.Println("Processing file")
}

func printNumber(number int) {
	fmt.Println("Deferred number:", number)
}

func main() {
	fmt.Println("----- Basic Defer -----")

	defer fmt.Println("This runs at the end of main")
	fmt.Println("This runs first")

	fmt.Println("----- Multiple Defers -----")

	defer fmt.Println("Third defer")
	defer fmt.Println("Second defer")
	defer fmt.Println("First defer")

	fmt.Println("Defer runs in LIFO order")

	fmt.Println("----- Defer For Cleanup -----")

	processFile()

	fmt.Println("----- Defer Argument Timing -----")

	number := 10
	defer printNumber(number)

	number = 20
	fmt.Println("Current number:", number)

	fmt.Println("Main function ending")
}
