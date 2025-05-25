package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This project explores practical LLM integration patterns in a customer support API scenario. Welcome:)")
	})

	api := r.Group("/api/support")
	{
		api.POST("/basic-llm-completion", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Basic LLM completion endpoint - simple Q&A with an LLM model",
			})
		})
	}

	port := "8080"
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
