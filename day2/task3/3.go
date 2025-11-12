package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Error: invalid amount of arguments")
		return
	}

	K, err := strconv.ParseUint(os.Args[1], 10, 64)

	if err != nil {
		fmt.Print("Error: Argument must have type uint")
		return
	}

	TickCounter := 0
	TimeCounter := K
	SignalChannel := make(chan os.Signal, 1)
	signal.Notify(SignalChannel, syscall.SIGINT, syscall.SIGTERM)

	go Ticker(TickCounter, TimeCounter, K)

	SigRes, ok := <-SignalChannel
	if ok {
		fmt.Println("Got signal:", SigRes)
		fmt.Println("Termination")
		return
	}
}

func Ticker(TickCounter int, TimeCounter, K uint64) {
	for {
		time.Sleep(time.Duration(K) * time.Second)
		TickCounter++
		fmt.Println("Tick", TickCounter, "since", TimeCounter)
		TimeCounter += K
	}
}