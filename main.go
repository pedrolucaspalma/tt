package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type CommandUsage struct {
	sessionId   string
	mainCommand string
	startTime   string
	closeTime   string
	fullCommand string
}

func main() {
	var commands []string = os.Args[1:]

	ttPtr, err := GetTt()
	if err != nil {
		log.Fatal(err)
	}
	if len(commands) > 0 {
		err := ttPtr.RegisterCommand(commands)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		ttPtr.GetReport()
	}
}

func (tt *Tt) GetReport() error {
	fmt.Println("Read execution")
	return nil
}

type Tt struct {
	TtConfigFilePath     string
	RegistrationFilePath string

	currentWorkingDirectoryPath string
}

func GetTt() (*Tt, error) {
	currentWorkingDirectoryPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	ttConfigFileName := "/tt_config.json"
	ttConfigFilePath := currentWorkingDirectoryPath + ttConfigFileName

	registrationFilePath := currentWorkingDirectoryPath + "/tt.csv"

	tt := Tt{
		TtConfigFilePath:            ttConfigFilePath,
		RegistrationFilePath:        registrationFilePath,
		currentWorkingDirectoryPath: currentWorkingDirectoryPath,
	}

	fmt.Println(tt)

	fmt.Println("Current Working Directory: ", currentWorkingDirectoryPath)
	return &tt, nil
}

func (tt *Tt) RegisterCommand(commands []string) error {
	sessionId, err := exec.Command("uuidgen").Output()
	if err != nil {
		return err
	}
	// removing \n from uuid
	sessionId = sessionId[:len(sessionId)-1]

	// creating  a string containing the full command (command + all arguments)
	fullCommandString := commands[0]
	for _, command := range commands[1:] {
		fullCommandString += " " + command
	}

	registrationFilePtr, err := tt.getOrCreateRegistrationFile()
	if err != nil {
		return err
	}
	defer registrationFilePtr.Close()

	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdout = os.Stdout

	startTime := getNow()

	err = cmd.Run()
	if err != nil {
		return err
	}

	closeTime := getNow()

	sessionData := CommandUsage{
		sessionId:   string(sessionId),
		mainCommand: commands[0],
		startTime:   startTime,
		closeTime:   closeTime,
		fullCommand: fullCommandString,
	}
	sessionString := generateString(sessionData)

	_, err = registrationFilePtr.WriteString(sessionString)
	if err != nil {
		return err
	}
	return nil
}

func (tt *Tt) getOrCreateRegistrationFile() (filePtr *os.File, err error) {
	filePath := tt.currentWorkingDirectoryPath + "/tt.csv"

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

func generateString(sessionData CommandUsage) string {
	return sessionData.sessionId + "," + sessionData.mainCommand + "," + sessionData.startTime + "," + sessionData.closeTime + "," + sessionData.fullCommand + ";" + "\n"
}

func getNow() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}
