package routes

import (
	"net/http"
	"strconv"

	"Watchdog/services"

	"github.com/gin-gonic/gin"
)

func SetupLogRoutes(router *gin.Engine) {
	logGroup := router.Group("/logs")
	{
		logGroup.GET("/type/:type", getLogsByType)
		logGroup.GET("/id/:id", getLogByID)
		logGroup.GET("/", getLogs)
		logGroup.GET("/history", getHistoricalLogs)

		logGroup.GET("/stored", getStoredLogs)
		logGroup.GET("/stored/id/:id", getStoredLogByID)
		logGroup.GET("/stored/type/:type", getStoredLogsByType)

		logGroup.POST("/stored/id/:id", postStoredLogsByID)
		logGroup.POST("/stored/type/:type", postLogsByType)

		logGroup.DELETE("/stored/id/:id", deleteStoredLogsByID)
		logGroup.DELETE("/stored/type/:type", deleteLogsByType)
	}
}

// History logs handlers /logs

// LOGS HANDLERS
// Get all logs handler
func getLogs(c *gin.Context) {
	events := services.LogsEvents()
	if len(events) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No logs found"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func getHistoricalLogs(c *gin.Context) {
	historyLogs := services.GetHistoricalLogs()
	if len(historyLogs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No historical logs found"})
		return
	}
	c.JSON(http.StatusOK, historyLogs)
}

// TYPE HANDLERS
// Get logs by type handler
func getLogsByType(c *gin.Context) {
	logType := c.Param("type")
	if logType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Log type cannot be empty"})
		return
	}
	if !services.IsValidLogType(logType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log type", "type": logType})
		return
	}
	result := services.GetLogsByType(logType)
	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No logs found for type", "type": logType})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ID HANDLERS
// Get log by ID handler
func getLogByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	logEntry, found := services.GetLogByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found", "id": id})
		return
	}
	c.JSON(http.StatusOK, logEntry)
}

// Stored logs handlers /logs/stored

// STORED LOGS HANDLERS
// Get all stored logs handler
func getStoredLogs(c *gin.Context) {
	logs := services.GetStoredLogs()
	if len(logs) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No stored logs found"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GET STORED ID HANDLERS
// Get stored log by ID handler
func getStoredLogByID(c *gin.Context) {
	logID := c.Param("id")
	id, err := strconv.Atoi(logID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	logEntry, found := services.GetStoredLogByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stored log not found", "id": id})
		return
	}
	c.JSON(http.StatusOK, logEntry)
}

// GET STORED TYPE HANDLERS
// Get stored logs by type handler
func getStoredLogsByType(c *gin.Context) {
	logType := c.Param("type")
	if logType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Log type cannot be empty"})
		return
	}
	if !services.IsValidLogType(logType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log type", "type": logType})
		return
	}
	result := services.GetStoredLogsByType(logType)
	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No stored logs found for type", "type": logType})
		return
	}
	c.JSON(http.StatusOK, result)
}

// POST and DELETE HANDLERS

// POST ID HANDLERS
// Post stored logs by ID handler
func postStoredLogsByID(c *gin.Context) {
	logID := c.Param("id")
	id, err := strconv.Atoi(logID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	services.PostLogByID(id)
	c.JSON(http.StatusOK, gin.H{"message": "Log posted successfully", "id": id})
}

// POST TYPE HANDLERS
// Post logs by type handler
func postLogsByType(c *gin.Context) {
	logType := c.Param("type")
	if !services.IsValidLogType(logType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log type", "type": logType})
		return
	}
	services.PostLogByType(logType)
	c.JSON(http.StatusOK, gin.H{"message": "Logs posted successfully", "type": logType})
}

// DELETE ID HANDLERS
// Delete stored logs by ID handler
func deleteStoredLogsByID(c *gin.Context) {
	logID := c.Param("id")
	id, err := strconv.Atoi(logID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	services.DeleteLogByID(id)
	c.JSON(http.StatusOK, gin.H{"message": "Log deleted successfully", "id": id})
}

// DELETE TYPE HANDLERS
// Delete logs by type handler
func deleteLogsByType(c *gin.Context) {
	logType := c.Param("type")
	if logType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Log type cannot be empty"})
		return
	}
	if !services.IsValidLogType(logType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log type", "type": logType})
		return
	}
	services.DeleteLogByType(logType)
	c.JSON(http.StatusOK, gin.H{"message": "Logs deleted successfully", "type": logType})
}
