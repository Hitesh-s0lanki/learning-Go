package main

import "fmt"

func main() {

	// for loop
	for i := 1; i < 10; i++ {
		fmt.Println(i)
	}

	// while loop
	k := 3
	for k != 0 {
		fmt.Println(k)
		k--
	}

	items := []string{"hitesh", "niraj", "kapil"}
	for _, value := range items {
		fmt.Println(value)
	}

}
