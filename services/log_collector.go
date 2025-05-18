package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Logs struct {
	TimeCreated      string `json:"timeCreated"` // Cambiado a string para compatibilidad con JSON de PowerShell
	ID               int    `json:"id"`
	LevelDisplayName string `json:"levelDisplayName"` //Type of log (INFO, WARN, ERROR, FATAL)
	Message          string `json:"message"`          //Message of the log
}

func parseWinDate(dateStr string) string {
	regularExpresion := regexp.MustCompile(`/Date\((\d+)\)/`)
	matches := regularExpresion.FindStringSubmatch(dateStr)
	if len(matches) == 2 {
		milliseconds, err := strconv.ParseInt(matches[1], 10, 64)
		if err == nil {
			t := time.Unix(0, milliseconds*int64(time.Millisecond))
			return t.Format(time.RFC3339)
		}
	}
	return dateStr
}

func LogsEvents() []Logs {
	cmd := exec.Command("powershell", "-Command", `Get-WinEvent -MaxEvents 10 | Select-Object TimeCreated, Id, LevelDisplayName, Message | ConvertTo-Json`)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error in powershell execution: %v\n", err)
		return nil
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

	for i := range events {
		events[i].TimeCreated = parseWinDate(events[i].TimeCreated)
	}

	return events
}

func GetLogByID(logs []Logs, id int) (Logs, bool) {
	for _, log := range logs {
		if log.ID == id {
			return log, true
		}
	}
	return Logs{}, false
}

func GetLogsByType(logs []Logs, logType string) []Logs {
	var result []Logs
	for _, log := range logs {
		if strings.EqualFold(log.LevelDisplayName, logType) {
			result = append(result, log)
		}
	}
	return result
}
