package main

import (
	"fmt"
)

type PatientInfo struct {
	Doctor string
	name [3]string
	date string
}

func main() {
	RunCommand(ReadCommand())
}

func ReadCommand() string {
	var Command string
	fmt.Scan(&Command)
	fmt.Print("\n")
	return Command
}

func RunCommand(Command string) {
	switch Command {
	case "Save":
		
	case "GetHistory":

	case "GetLastVisit":

	case "Exit", "exit":
	}
}