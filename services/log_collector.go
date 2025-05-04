package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type Logs struct {
	ID               int       `json:"id"`
	TimeCreated      time.Time `json:"timeCreated"`      //Time at which the log was generated
	LevelDisplayName string    `json:"levelDisplayName"` //Type of log (INFO, WARN, ERROR, FATAL)
	Message          string    `json:"message"`          //Message of the log
}

func logsEvents() []Logs {
	cmd := exec.Command("powershell", "-Command", `Get-WinEvent -LogName System -MaxEvents 10 | ConvertTo-Json`)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error in powershell execution: %v\n", err)
		return nil
	}

	var events []Logs
	if err := json.Unmarshal(out.Bytes(), &events); err != nil {

	}
}
