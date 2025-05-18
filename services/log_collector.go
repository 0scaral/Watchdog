package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type Logs struct {
	TimeCreated      time.Time `json:"timeCreated"` //Time at which the log was generated
	ID               int       `json:"id"`
	LevelDisplayName string    `json:"levelDisplayName"` //Type of log (INFO, WARN, ERROR, FATAL)
	Message          string    `json:"message"`          //Message of the log
}

func LogsEvents() ([]Logs, map[int]Logs) {
	cmd := exec.Command("powershell", "-Command", `Get-WinEvent -LogName System -MaxEvents 10 | ConvertTo-Json`)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error in powershell execution: %v\n", err)
		return nil, nil
	}

	var events []Logs
	if err := json.Unmarshal(out.Bytes(), &events); err != nil {
		var event Logs
		if err := json.Unmarshal(out.Bytes(), &event); err == nil {
			events = append(events, event)
		} else {
			fmt.Printf("Error parse JSON: %v", err)
		}
	}

	logMap := make(map[int]Logs)
	for _, log := range events {
		logMap[log.ID] = log
	}

	return events, logMap
}

func GetLogByID(logMap map[int]Logs, id int) (Logs, bool) {
	log, found := logMap[id]
	return log, found
}

func GetLogsByType(logMap map[int]Logs, logType string) []Logs {
	var result []Logs
	for _, log := range logMap {
		if log.LevelDisplayName == logType {
			result = append(result, log)
		}
	}
	return result
}
