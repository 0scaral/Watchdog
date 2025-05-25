package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Logs struct {
	TimeCreated      string `json:"timeCreated"`
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

func fetchLogs() []Logs {
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

var (
	storedLogs      []Logs
	storedLogsMutex sync.RWMutex
)

// addLogsToStored agrega logs a storedLogs si no están ya presentes.
func addLogsToStored(logs []Logs) {
	storedLogsMutex.Lock()
	defer storedLogsMutex.Unlock()
	existing := make(map[int]struct{})
	for _, l := range storedLogs {
		existing[l.ID] = struct{}{}
	}
	for _, log := range logs {
		if _, found := existing[log.ID]; !found {
			storedLogs = append(storedLogs, log)
			existing[log.ID] = struct{}{}
		}
	}
}

// LogsEvents obtiene los logs más recientes y los almacena en storedLogs.
func LogsEvents() []Logs {
	logs := fetchLogs()
	for _, log := range logs {
		switch strings.ToLower(log.LevelDisplayName) {
		case "error", "critical", "warning":
			msg := fmt.Sprintf("Log Alert, a suspicious log has been detected.\nID: %d\nType: %s\nMessage: %s", log.ID, log.LevelDisplayName, log.Message)
			sendAlerts(msg)
		}
	}
	addLogsToStored(logs)
	return logs
}

// Consultas sobre storedLogs en vez de logs pasados por parámetro.
func GetLogByID(id int) (Logs, bool) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range storedLogs {
		if log.ID == id {
			return log, true
		}
	}
	return Logs{}, false
}

func GetLogsByType(logType string) []Logs {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	var result []Logs
	for _, log := range storedLogs {
		if strings.EqualFold(log.LevelDisplayName, logType) {
			result = append(result, log)
		}
	}
	return result
}
