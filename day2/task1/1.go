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

			fmt.Println("goroutine", id, "slept", sleeptime, "ms")
		}(i)
	}

	wg.Wait()

	sortarray := make([]int, 0, len(m))
	for k := range m {
		sortarray = append(sortarray, k)
	}

	sort.Slice(sortarray, func(i, j int) bool {
		return m[sortarray[i]] < m[sortarray[j]]
	})

	fmt.Println("All goroutines finished. Collected results:")
	for _, id := range sortarray {
		fmt.Print(id, m[id], "\n")
	}
}
