package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
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
	historyLogs     []Logs
	storedLogs      []Logs
	storedLogsMutex sync.RWMutex
)

func addLogsToHistory(logs []Logs) {
	storedLogsMutex.Lock()
	defer storedLogsMutex.Unlock()
	existing := make(map[int]struct{})
	for _, l := range historyLogs {
		existing[l.ID] = struct{}{}
	}
	for _, log := range logs {
		if _, found := existing[log.ID]; !found {
			historyLogs = append(historyLogs, log)
			existing[log.ID] = struct{}{}
		}
	}
}

func LogsEvents() []Logs {
	addLogsToHistory(fetchLogs())
	return slices.Clone(fetchLogs())
}

func GetLogByID(id int) (Logs, bool) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range historyLogs {
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
	for _, log := range historyLogs {
		if strings.EqualFold(log.LevelDisplayName, logType) {
			result = append(result, log)
		}
	}
	return result
}

func PostLogByID(id int) error {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range historyLogs {
		if log.ID == id {
			storedLogs = append(storedLogs, log)
		} else {
			return fmt.Errorf("log with ID %d not found", id)
		}
	}
	return nil
}

func PostLogByType(logType string) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range historyLogs {
		if strings.EqualFold(log.LevelDisplayName, logType) {
			storedLogs = append(storedLogs, log)
		}
	}
}

func DeleteLogByID(id int) error {
	storedLogsMutex.Lock()
	defer storedLogsMutex.Unlock()
	for i, log := range storedLogs {
		if log.ID == id {
			storedLogs = append(storedLogs[:i], storedLogs[i+1:]...)
		} else {
			return fmt.Errorf("log with ID %d not found", id)
		}
	}
	return nil
}

func DeleteLogByType(logType string) {
	storedLogsMutex.Lock()
	defer storedLogsMutex.Unlock()
	var newLogs []Logs
	for _, log := range storedLogs {
		if !strings.EqualFold(log.LevelDisplayName, logType) {
			newLogs = append(newLogs, log)
		}
	}
	storedLogs = newLogs
}

func GetStoredLogs() []Logs {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	return slices.Clone(storedLogs)
}

func GetStoredLogsByType(logType string) []Logs {
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

func GetStoredLogByID(id int) (Logs, bool) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range storedLogs {
		if log.ID == id {
			return log, true
		}
	}
	return Logs{}, false
}

var validLogTypes = []string{"Information", "Warning", "Error", "Critical", "Verbose"}

// IsValidLogType checks if the provided log type is valid (case-insensitive)
func IsValidLogType(logType string) bool {
	for _, t := range validLogTypes {
		if strings.EqualFold(t, logType) {
			return true
		}
	}
	return false
}
