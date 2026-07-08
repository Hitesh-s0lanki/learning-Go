package main

import "fmt"

func sum(numbers ...int) int {
	sum := 0
	for _, val := range numbers {
		sum += val
	}

	return sum
}

func main() {

	fmt.Println(sum(1, 2, 3, 4, 5))

}
