package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	FirstBuffer := readstr()
	FirstNumbers, err := AppendToSlice(FirstBuffer)
	if err != nil {
		fmt.Print("Error: invalid input")
		return
	}

	SecondBuffer := readstr()
	SecondNumbers, err := AppendToSlice(SecondBuffer)
	if err != nil {
		fmt.Print("Error: invalid input")
		return
	}
	
	SimilarNumbers := FindIntersection(FirstNumbers, SecondNumbers)
	if SimilarNumbers == nil {
		fmt.Print("Empty intersection")
	}

	PrintResult(SimilarNumbers)
}

func readstr() []string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	buffer := strings.Fields(line)
	return buffer
}

func AppendToSlice(numbers []string) ([]int, error) {
	var NumbersSlice []int
	for _, i := range numbers {
		i, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		NumbersSlice = append(NumbersSlice, i)
	}
	return NumbersSlice, nil
}

func FindIntersection(FirstNumbers, SecondNumbers []int) []int {
	var SimilarNumbers []int
	for _, i := range FirstNumbers {
		for _, j := range SecondNumbers {
			if i == j {
				SimilarNumbers = append(SimilarNumbers, j)
				break
			}
		}
	}
	return SimilarNumbers
}

func PrintResult(SimilarNumbers []int) {
	for i, number := range SimilarNumbers {
		if i == len(SimilarNumbers) - 1 {
			fmt.Print(number)
		} else {
			fmt.Print(number, " ")
		}
	} 
}