package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Error: invalid amount of arguments")
		return
	}

	K, err := strconv.ParseUint(os.Args[1], 10, 64)

	if err != nil {
		fmt.Print("Error: Argument must have type int")
		return
	}


}