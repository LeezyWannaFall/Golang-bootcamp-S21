package main

import (
	"fmt"
)

type PatientInfo struct {
	doctor string
	name string
	date string
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