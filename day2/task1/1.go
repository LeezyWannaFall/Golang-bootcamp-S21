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

	if !CheckInt(N, M) {
		return
	}


	m := make(map[int]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	CreateRoutines(&wg, &mu, N, M, m)
	wg.Wait() // блокаем main горутину (ждем другие горутины)
	SortResult(m)
}

func SortResult(m map[int]int) {
	var results []Result
	for id, sleeptime := range m {
		results = append(results, Result{id, sleeptime})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].sleeptime > results[j].sleeptime
	})

	fmt.Println("All goroutines finished. Collected results:")
	for _, i := range results {
		fmt.Println(i.id, i.sleeptime)
	}
}

func CreateRoutines(wg *sync.WaitGroup, mu *sync.Mutex, N, M int, m map[int]int) {
	for i := 1; i < N + 1; i++ {
		wg.Add(1) // добавляем горутину (+1)
		go func(id int) {
			defer wg.Done() // убавляем горутину (-1) 
			sleeptime := rand.Intn(M)
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)

			mu.Lock()
			m[id] = sleeptime
			mu.Unlock()

			fmt.Print("Goroutine #", id, sleeptime, " ms\n")
		}(i)
	}
}

func CheckInt(N, M int) bool {
	if M <= 0 || N <= 0 {
		fmt.Print("Numbers must be positive")
		return false
	}
	return true
}