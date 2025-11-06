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

type Result struct {
	id, sleeptime int
}

func main() {
	if len(os.Args) != 3 {
		fmt.Print("Error: invalid amount of arguments")
		return
	}

	N, err := strconv.Atoi(os.Args[1])
	
	if err != nil {
		fmt.Print("Error: arguments must have type int")
		return
	}

	M, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Print("Error: arguments must have type int")
		return
	}
	
}