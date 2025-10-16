package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type WordCount struct {
    word  string
    count int
}

func main() {
	var m map[string]int
	var result []WordCount
	var K int

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	words := strings.Fields(line)

	fmt.Scan(&K)
	m = make(map[string]int)
	for i := 0; i < len(words); i++ {
		m[words[i]]++
	}

	for k, v := range m {
		result = append(result, WordCount{k, v})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].count == result[j].count {
			return result[i].word < result[j].word
		} else {
			return result[i].count > result[j].count
		}
	})

	if K > len(result) {
		for i := 0; i < len(result); i++ {
			if i == len(result) - 1 {
				fmt.Print(result[i].word)	
			} else {
				fmt.Print(result[i].word, " ")
			}
		}
	} else {
		for i := 0; i < K; i++ {
			if i == K - 1 {
				fmt.Print(result[i].word)	
			} else {
				fmt.Print(result[i].word, " ")
			}
		}
	}
}