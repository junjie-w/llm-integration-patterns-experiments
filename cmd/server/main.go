package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/basic_llm_completion"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/api/handlers"
	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	basicLLMCompletionService := basic_llm_completion.NewService(cfg)

	basicLLMCompletionHandler := handlers.NewCompletionHandler(basicLLMCompletionService)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This project explores practical LLM integration patterns in a customer support API scenario. Welcome:)")
	})

	api := r.Group("/api/support")
	{
		api.POST("/basic-llm-completion", basicLLMCompletionHandler.HandleBasicCompletion)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
