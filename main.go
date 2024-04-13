package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	var commands []string = os.Args[1:]

	if len(commands) < 1 {
		log.Fatal("No command found")
	}

	sessionId, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	// removing \n from uuid
	sessionId = sessionId[:len(sessionId)-1]

	// creating  a string containing the full command (command + all arguments)
	fullCommandString := commands[0]
	for _, command := range commands[1:] {
		fullCommandString += " " + command
	}

	registrationFilePtr, err := getOrCreateRegistrationFile()
	if err != nil {
		log.Fatal(err)
	}
	defer registrationFilePtr.Close()

	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdout = os.Stdout

	startTime := getNow()

	cmd.Run()

	closeTime := getNow()

	sessionString := generateString(string(sessionId), commands[0], startTime, closeTime, fullCommandString)

	_, err = registrationFilePtr.WriteString(sessionString)
	if err != nil {
		log.Fatal(err)
	}
}

func getNow() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

func getOrCreateRegistrationFile() (filePtr *os.File, err error) {
	filePath := "/home/pedro/tt.csv"

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {
		file, err = os.Create(filePath)
		file.WriteString("sessionId,command,startTime,endTime,fullCommandString;\n")
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func generateString(sessionId string, command string, start string, end string, fullCommandString string) string {
	return sessionId + "," + command + "," + start + "," + end + "," + fullCommandString + ";" + "\n"
}
