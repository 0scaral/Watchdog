package routes

import (
	"net/http"
	"strconv"

	"Watchdog/services"

	"github.com/gin-gonic/gin"
)

func SetupLogRoutes(router *gin.Engine) {
	logGroup := router.Group("/api/logs")
	{
		logGroup.GET("/", getLogs)
		logGroup.GET("/type/:type", getLogsByType)
		logGroup.GET("/id/:id", getLogByID)
	}
}

func getLogs(c *gin.Context) {
	events, _ := services.LogsEvents()
	c.JSON(http.StatusOK, events)
}

func getLogsByType(c *gin.Context) {
	_, logMap := services.LogsEvents()
	logType := c.Param("type")
	result := services.GetLogsByType(logMap, logType)
	c.JSON(http.StatusOK, result)
}

func getLogByID(c *gin.Context) {
	_, logMap := services.LogsEvents()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	log, found := services.GetLogByID(logMap, id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}
	c.JSON(http.StatusOK, log)
}
