package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type PatientInfo struct {
	doctor string
	date string
}

func main() {
	RunCommand()
}

func RunCommand() {
	m := make(map[string][]PatientInfo)

	for {
		var Command string
		fmt.Print("> ")
		fmt.Scan(&Command)
		
		switch Command {
		case "Save", "save":
			var doctor, date, name string
			ReadName()
			fmt.Scan(&doctor)
			fmt.Scan(&date)

			m[name] = append(m[name], PatientInfo{doctor: doctor, date: date})
		case "GetHistory", "gethistory":
			var name string
			ReadName()

			patient, ok := m[name]
			if !ok {
				fmt.Println("patient not found")
				return
			}

			for _, i := range patient {
				fmt .Println("Doctor: ", i.doctor, "| Date: ", i.date)
			}
		case "GetLastVisit", "getlastvisit":
			
		case "Exit", "exit":
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}

func ReadName() string {
	reader := bufio.NewReader(os.Stdin)
	line,_ := reader.ReadString('\n')
	name := strings.TrimSpace(line)
	return name
}