package main

import (
	"Watchdog/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Static("/static", "./static")

	routes.SetupLogRoutes(router)

	log.Println("Starting server on :8080")
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
