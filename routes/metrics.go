package routes

import (
	"net/http"
	"time"

	"Watchdog/services"

	"github.com/gin-gonic/gin"
)

func SetupMetricsRoutes(router *gin.Engine) {
	metricsGroup := router.Group("/metrics")
	{
		metricsGroup.GET("/cpu", getCPUUsage)
		metricsGroup.GET("/ram", getRAMUsage)
		metricsGroup.GET("/disk", getDiskUsage)
		metricsGroup.GET("/current", getCurrentMetrics)
		metricsGroup.GET("/average/:minutes", getAverageMetrics)
	}
}

func getCPUUsage(c *gin.Context) {
	metric := services.GetCurrentMetric()
	c.JSON(http.StatusOK, gin.H{"cpu_usage": metric.CPUUsage})
}

func getRAMUsage(c *gin.Context) {
	metric := services.GetCurrentMetric()
	c.JSON(http.StatusOK, gin.H{"ram_usage": metric.RAMUsage})
}

func getDiskUsage(c *gin.Context) {
	metric := services.GetCurrentMetric()
	c.JSON(http.StatusOK, gin.H{"disk_usage": metric.DiskUsage})
}

func getCurrentMetrics(c *gin.Context) {
	metric := services.GetCurrentMetric()
	c.JSON(http.StatusOK, metric)
}

func getAverageMetrics(c *gin.Context) {
	minutesStr := c.Param("minutes")
	minutes, err := time.ParseDuration(minutesStr + "m")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minutes parameter"})
		return
	}
	metric := services.GetAverageMetric(minutes)
	c.JSON(http.StatusOK, metric)
}
