package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize Kubernetes client
	clientset, err := NewKubeClient()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	// Middleware to get userID from header
	r.Use(func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing X-User-ID header"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	})

	r.GET("/pods", GetPodsHandler(clientset))
	r.GET("/search", SearchPodsHandler(clientset))

	r.Run(":8080")
}
