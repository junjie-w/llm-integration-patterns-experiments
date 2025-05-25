package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/knowledge_rag"
)

type KnowledgeRagHandler struct {
	service *knowledge_rag.Service
}

func NewKnowledgeRagHandler(service *knowledge_rag.Service) *KnowledgeRagHandler {
	return &KnowledgeRagHandler{
		service: service,
	}
}

func (h *KnowledgeRagHandler) HandleKnowledgeRagCompletion(c *gin.Context) {
	var req knowledge_rag.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message cannot be empty"})
		return
	}

	resp, err := h.service.GetCompletion(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get knowledge completion"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
