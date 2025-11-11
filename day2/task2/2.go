package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Print("Error: invalid amount of arguments")
		return
	}

	K, err := strconv.Atoi(os.Args[1])
	
	if err != nil {
		fmt.Print("Error: arguments must have type int")
		return
	}

	N, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Print("Error: arguments must have type int")
		return
	}

	if !CheckInt(N, K) {
		return
	}

	pow := Pow(Generator(K, N))
	for num := range pow {
		fmt.Println(num)
	}
	
}

func Generator(K, N int) <-chan int {
	firstchan := make(chan int)
	go func() {
		for i := K; i <= N; i++ {
			firstchan <- i
		}
		close(firstchan)
	}()
	return firstchan
}

func Pow(firstchan <-chan int) <-chan int {
	secondchan := make(chan int)
	go func() {
		for val := range firstchan {
			secondchan <- val * val
		}
		close(secondchan)
	}()
	return secondchan
}


func CheckInt(N, K int) bool {
	if K > N {
		fmt.Print("Error: First argument must be lower than second")
		return false
	}
	return true
}