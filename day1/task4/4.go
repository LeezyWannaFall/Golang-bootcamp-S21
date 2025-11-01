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

type PatientNotFoundError struct {
	Message string
}

func (e *PatientNotFoundError) Error() string {
	return e.Message
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
			err := GetHistory(m, name)

			if err != nil {
				fmt.Print(err)
			}
		case "GetLastVisit", "getlastvisit":
			var doctor string
			name := ReadName()
			fmt.Scanln(&doctor)
			ValidDoctor := true

			err, ValidDoctor := GetLastVisit(m, name, doctor, ValidDoctor)

			if err != nil {
				fmt.Print(err)
			} else {
				if !ValidDoctor {
					fmt.Println("invalid doctor")
				}
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

func GetHistory(m map[string][]PatientInfo, name string) error {
	visits, ok := m[name]
	if !ok {
		return &PatientNotFoundError{Message: "patient not found\n"}
	} else {
		for _, v := range visits {
			fmt.Print("Doctor: ", v.doctor, " Date: ", v.date, "\n")
		}
	}
	return nil
}

func GetLastVisit(m map[string][]PatientInfo, name, doctor string, ValidDoctor bool) (error, bool) {
	visits, ok := m[name]
	if !ok {
		return &PatientNotFoundError{Message: "patient not found\n"}, false
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
	return nil, ValidDoctor
}