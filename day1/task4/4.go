package main

import (
	"fmt"
)

type PatientInfo struct {
	doctor string
	date string
}

type PatientName struct {
	name string
	PatientInfo struct{}
}

func main() {
	RunCommand()
}

func RunCommand() {
	m := make(map[string]PatientInfo)

	for {
		var Command string
		fmt.Print("> ")
		fmt.Scan(&Command)
		
		switch Command {
		case "Save", "save":
			var doctor, name, date string
			fmt.Scan(&name, &doctor, date)
			// ...
		case "GetHistory", "gethistory":
			var name string
		case "GetLastVisit", "getlastvisit":
			
		case "Exit", "exit":
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}