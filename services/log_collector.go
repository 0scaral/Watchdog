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

	models "Watchdog/models"
)

var (
	historyLogs     []models.Logs
	storedLogs      []models.Logs
	storedLogsMutex sync.RWMutex
	alertedLogs     = make(map[string]struct{})
	alertedLogMutex sync.Mutex
)

// validLogTypes contains the valid log types for filtering
var validLogTypes = []string{"Information", "Warning", "Error", "Critical", "Verbose"}

func IsValidLogType(logType string) bool {
	for _, t := range validLogTypes {
		if strings.EqualFold(t, logType) {
			return true
		}
	}
	return false
}

// parseWinDate converts a Windows date string to RFC3339 format
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

// Obtain the lastest logs from Windows Event Viewer
func fetchLogs() []models.Logs {
	cmd := exec.Command("powershell", "-Command", `Get-WinEvent -LogName 'Application','System','Security' -MaxEvents 10 -ErrorAction SilentlyContinue | Select-Object TimeCreated, Id, LevelDisplayName, Message | ConvertTo-Json -Depth 5`)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error in powershell execution: %v\n", err)
		return nil
	}

	var events []models.Logs
	if err := json.Unmarshal(out.Bytes(), &events); err != nil {
		var event models.Logs
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

// addLogsToHistory adds new logs to the history, avoiding duplicates
func addLogsToHistory(logs []models.Logs) {
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

func logUniqueKey(log models.Logs) string {
	return fmt.Sprintf("%d_%s", log.ID, log.TimeCreated)
}

// HISTORICAL LOGS HANDLERS
// Function to handle log events and send alerts
func LogsEvents() []models.Logs {
	logs := fetchLogs()
	for _, log := range logs {
		shouldAlert := false
		switch strings.ToLower(log.LevelDisplayName) {
		case "error", "critical", "warning":
			key := logUniqueKey(log)
			alertedLogMutex.Lock()
			_, alreadyAlerted := alertedLogs[key]
			if !alreadyAlerted {
				alertedLogs[key] = struct{}{}
				shouldAlert = true
			}
			alertedLogMutex.Unlock()
			if shouldAlert {
				msg := fmt.Sprintf("Log Alert, a suspicious log has been detected.\nID: %d\nType: %s\nMessage: %s", log.ID, log.LevelDisplayName, log.Message)
				SendAlertMail(msg)
				SendAlertTelegram(msg)
			}
		}
	}
	addLogsToHistory(logs)
	return logs
}

// GetLogs returns the history of logs
func GetLogByID(id int) (models.Logs, bool) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range historyLogs {
		if log.ID == id {
			return log, true
		}
	}
	return models.Logs{}, false
}

// GetLogsByType returns logs filtered by type
func GetLogsByType(logType string) []models.Logs {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	var result []models.Logs
	for _, log := range historyLogs {
		if strings.EqualFold(log.LevelDisplayName, logType) {
			result = append(result, log)
		}
	}
	return result
}

// GetHistoricalLogs returns the history of logs
func GetHistoricalLogs() []models.Logs {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	return slices.Clone(historyLogs)
}

// STORED LOGS HANDLERS
// GETTER
// GetStoredLogs returns all stored logs
func GetStoredLogs() []models.Logs {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	return slices.Clone(storedLogs)
}

// GetStoredLogsByType returns stored logs filtered by type
func GetStoredLogsByType(logType string) []models.Logs {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	var result []models.Logs
	for _, log := range storedLogs {
		if strings.EqualFold(log.LevelDisplayName, logType) {
			result = append(result, log)
		}
	}
	return result
}

// GetStoredLogByID returns a stored log by its ID
func GetStoredLogByID(id int) (models.Logs, bool) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range storedLogs {
		if log.ID == id {
			return log, true
		}
	}
	return models.Logs{}, false
}

// POST
// PostLogByID adds a log to the stored logs by its ID
func PostLogByID(id int) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range historyLogs {
		if log.ID == id {
			storedLogs = append(storedLogs, log)
		}
	}
}

// PostLogByType adds logs of a specific type to the stored logs
func PostLogByType(logType string) {
	storedLogsMutex.RLock()
	defer storedLogsMutex.RUnlock()
	for _, log := range historyLogs {
		if strings.EqualFold(log.LevelDisplayName, logType) {
			storedLogs = append(storedLogs, log)
		}
	}
}

// DELETE
// DeleteLogByID removes a log from the stored logs by its ID
func DeleteLogByID(id int) {
	storedLogsMutex.Lock()
	defer storedLogsMutex.Unlock()
	for i, log := range storedLogs {
		if log.ID == id {
			storedLogs = append(storedLogs[:i], storedLogs[i+1:]...)
		}
	}
}

// DeleteLogByType removes logs of a specific type from the stored logs
func DeleteLogByType(logType string) {
	storedLogsMutex.Lock()
	defer storedLogsMutex.Unlock()
	var newLogs []models.Logs
	for _, log := range storedLogs {
		if !strings.EqualFold(log.LevelDisplayName, logType) {
			newLogs = append(newLogs, log)
		}
	}
	storedLogs = newLogs
}

func StartLogCollection() {
	go func() {
		for {
			LogsEvents()
			time.Sleep(10 * time.Second)
		}
	}()
}
