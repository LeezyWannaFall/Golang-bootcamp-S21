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
			name := ReadName()
			var doctor, date string
			fmt.Scan(&doctor)
			fmt.Scan(&date)

			m[name] = append(m[name], PatientInfo{doctor: doctor, date: date})
		case "GetHistory", "gethistory":
			name := ReadName()
			
			visits, ok := m[name]
			if !ok {
				fmt.Println("patient not found")
			} else {
				for _, v := range visits {
					fmt.Print("Doctor: ", v.doctor, "\n", "Date: ", v.date, "\n")
				}
			}
		case "GetLastVisit", "getlastvisit":
			var doctor string
			name := ReadName()
			fmt.Scanln(&doctor)
			ValidDoctor := true
			visits, ok := m[name]
			if !ok {
				fmt.Print("patient not found\n")
			} else {
				for _, v := range visits {
					if doctor == v.doctor {
						fmt.Println("Date: ", v.date)
						ValidDoctor = true
					} else {
						ValidDoctor = false
					}
				}
			}

			if !ValidDoctor {
				fmt.Println("invalid doctor")
			}
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