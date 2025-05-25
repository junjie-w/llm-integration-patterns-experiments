package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/basic_llm_completion"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/embeddings"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/function_calling"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/knowledge_rag"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/multi_agent"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/reasoning_agent"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/tool"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/api/handlers"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/store/document"
	"github.com/junjie-w/llm-integration-patterns-experiments/pkg/config"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	_ = rng

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	docRepo := document.NewRepository()
	document.SeedDocuments(docRepo)

	toolRegistry := tool.NewRegistry()
	tool.RegisterSupportTools(toolRegistry)

	basicLLMCompletionService := basic_llm_completion.NewService(cfg)
	embeddingService := embeddings.NewService(cfg)
	knowledgeService := knowledge_rag.NewService(cfg, docRepo, embeddingService)
	functionCallingService := function_calling.NewService(cfg, toolRegistry)
	reasoningAgentService := reasoning_agent.NewService(cfg, toolRegistry)
	multiAgentService := multi_agent.NewService(cfg)

	basicLLMCompletionHandler := handlers.NewBasicLLMCompletionHandler(basicLLMCompletionService)
	knowledgeHandler := handlers.NewKnowledgeRagHandler(knowledgeService)
	functionCallingHandler := handlers.NewFunctionCallingHandler(functionCallingService)
 	reasoningAgentHandler := handlers.NewReasoningAgentHandler(reasoningAgentService)
  multiAgentHandler := handlers.NewMultiAgentHandler(multiAgentService)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This project explores practical AI development patterns in a customer support API scenario. Welcome:)")
	})

	api := r.Group("/api/support")
	{
		api.POST("/basic-llm-completion", basicLLMCompletionHandler.HandleBasicLLMCompletion)
		api.POST("/knowledge-rag", knowledgeHandler.HandleKnowledgeRagCompletion)
		api.POST("/function-calling", functionCallingHandler.HandleFunctionCallingCompletion)
		api.POST("/reasoning-agent", reasoningAgentHandler.HandleReasoningAgentExecution)
		api.POST("/multi-agent", multiAgentHandler.HandleMultiAgentProcess)
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
