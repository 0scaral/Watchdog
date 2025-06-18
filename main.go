package main

import (
	"Watchdog/routes"
	"Watchdog/services"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	services.StartLogCollection()     // Inicia la recolección de logs
	services.StartMetricsCollection() // Inicia la recolección de métricas

	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	router.Static("/static", "./static")

	routes.SetupLogRoutes(router)
	routes.SetupMetricsRoutes(router)

	log.Println("Starting server on :8080")
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
