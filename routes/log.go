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
	}
}

func getLogs(c *gin.Context) {
	events := services.LogsEvents()
	if len(events) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No logs found"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func getLogsByType(c *gin.Context) {
	logType := c.Param("type")
	result := services.GetLogsByType(logType)
	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No logs found for type", "type": logType})
		return
	}
	c.JSON(http.StatusOK, result)
}

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
