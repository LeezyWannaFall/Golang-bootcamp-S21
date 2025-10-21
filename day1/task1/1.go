package main

import (
	"fmt"
)

func main() {
	var x, y float64
	var operation string
	scan(&x, &y, &operation)
	calculate(x, y, operation)
}

func scan(x, y *float64, operation *string) {
	for {
		fmt.Print("Input right operand: ")
		if _, err := fmt.Scan(x); err == nil {
			break
		}
		fmt.Print("Invalid input. Please enter a number.\n")
		clearInputBuffer()
	}

	ValidOperators := []string{"+", "-", "*", "/"}
	for {
		fmt.Print("Input operation (+, -, *, /): ")
		fmt.Scan(operation)

		isValid := false
		for _, v := range ValidOperators {
			if *operation == v {
				isValid = true
				break
			}
		}

		if isValid {
			break
		}

		fmt.Print("Invalid operation. Please enter valid operations (+, -, *, /).\n")
	}

	for {
		fmt.Print("Input left operand: ")
		if _, err := fmt.Scan(y); err == nil {
			break
		}
		fmt.Print("Invalid input. Please enter a number.\n")
		clearInputBuffer()
	}
}

func calculate(x, y float64, operation string) {
	switch operation {
	case "+":
		fmt.Printf("Result: %.3f", x + y)
	case "-":
		fmt.Printf("Result: %.3f", x - y)
	case "*":
		fmt.Printf("Result: %.3f", x * y)
	case "/":
		if y == 0 {
			fmt.Print("Error: division by zero")
		} else {
			fmt.Printf("Result: %.3f", x / y)
		}
	}
}

func clearInputBuffer() {
    var discard string
    fmt.Scanln(&discard)
}

