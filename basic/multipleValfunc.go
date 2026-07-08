package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("Divide by zero not possible")
	}

	return a / b, nil
}

func main() {

	value, error := divide(1, 0)

	if error != nil {
		fmt.Println(error)
	} else {
		fmt.Println(value)
	}

}
