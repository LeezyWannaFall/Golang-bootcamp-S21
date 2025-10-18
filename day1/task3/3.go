package main

import (
	"fmt"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	words := strings.Fields(line)
}