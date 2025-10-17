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
	var result []WordCount
	var K int

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	words := strings.Fields(line)
	fmt.Scan(&K)

	MapAppend(words, &result)
	SortSlice(result)
	PrintResult(K, result)
}

func MapAppend(words []string, result *[]WordCount) {
	m := make(map[string]int)
	
	for i := 0; i < len(words); i++ {
		m[words[i]]++
	}

	for k, v := range m {
		*result = append(*result, WordCount{k, v})
	}
}

func SortSlice(result []WordCount) {
	sort.Slice(result, func(i, j int) bool {
		if result[i].count == result[j].count {
			return result[i].word < result[j].word
		} else {
			return result[i].count > result[j].count
		}
	})
}

func PrintResult(K int, result []WordCount) {
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
