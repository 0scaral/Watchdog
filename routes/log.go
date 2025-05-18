package routes

import (
	"log"
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
		log.Println("No logs found")
		c.JSON(http.StatusOK, gin.H{"message": "No logs found"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func getLogsByType(c *gin.Context) {
	events := services.LogsEvents()
	logType := c.Param("type")
	result := services.GetLogsByType(events, logType)
	if len(result) == 0 {
		log.Printf("No logs found for type: %s", logType)
		c.JSON(http.StatusNotFound, gin.H{"error": "No logs found for type", "type": logType})
		return
	}
	c.JSON(http.StatusOK, result)
}

func getLogByID(c *gin.Context) {
	events := services.LogsEvents()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	logEntry, found := services.GetLogByID(events, id)
	if !found {
		log.Printf("Log not found for ID: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found", "id": id})
		return
	}
	c.JSON(http.StatusOK, logEntry)
}
