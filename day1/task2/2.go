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
	var K int

	words := readstr()
	if len(words) == 0 {
		fmt.Print("\n")
		return 
	}

	MostFamousWords(K, words)
}

func readstr() []string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	words := strings.Fields(line)
	return words
}

func MapAppend(words []string, result *[]WordCount) []WordCount {
	m := make(map[string]int)
	
	for i := 0; i < len(words); i++ {
		m[words[i]]++
	}

	for k, v := range m {
		*result = append(*result, WordCount{k, v})
	}

	return *result
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

func WordCountToString(K int, result []WordCount) string {
	var words []string
	if K > len(result) {
		for i := 0; i < len(result); i++ {
			words = append(words, result[i].word)
		}
	} else {
		for i := 0; i < K; i++ {
			words = append(words, result[i].word)
		}
	}
	return strings.Join(words, " ")
}

func MostFamousWords(K int, words []string) string {
	var result []WordCount
	fmt.Scan(&K)

	SortSlice(MapAppend(words, &result))
	PrintResult(K, result)
	return WordCountToString(K, result)
}