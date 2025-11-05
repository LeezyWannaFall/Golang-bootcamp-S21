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
	N, _ := strconv.Atoi(os.Args[1])
	M, _ := strconv.Atoi(os.Args[2])
	m := make(map[int]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 1; i < N + 1; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sleeptime := rand.Intn(M)
			time.Sleep(time.Duration(sleeptime) * time.Millisecond)

			mu.Lock()
			m[id] = sleeptime
			mu.Unlock()

			fmt.Print("Goroutine #", id, sleeptime, " ms\n")
		}(i)
	}

	wg.Wait()


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
