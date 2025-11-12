package main

import (
	"log"
	"os"

	"codeheist/game"
	"codeheist/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize game engine
	gameEngine := game.NewEngine()

	// Start cleanup goroutine for expired sessions
	go gameEngine.CleanupSessions()

	// Initialize WebSocket handler
	wsHandler := websocket.NewHandler(gameEngine)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	router.GET("/ws", wsHandler.HandleWebSocket)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "codeheist"})
	})

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 CodeHeist server starting on port %s", port)
	log.Printf("💻 Web terminal available at http://localhost:%s", port)
	log.Printf("🔌 WebSocket endpoint: ws://localhost:%s/ws", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
