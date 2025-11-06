package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
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



}

func CheckInt(N, K int) bool {
	if K <= 0 || N <= 0 {
		fmt.Print("Numbers must be positive")
		return false
	}
	return true
}