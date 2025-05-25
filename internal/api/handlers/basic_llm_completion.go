package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/basic_llm_completion"
)

type BaiscLLMCompletionHandler struct {
	service *basic_llm_completion.Service
}

func NewBasicLLMCompletionHandler(service *basic_llm_completion.Service) *BaiscLLMCompletionHandler {
	return &BaiscLLMCompletionHandler{
		service: service,
	}
}

func (h *BaiscLLMCompletionHandler) HandleBasicLLMCompletion(c *gin.Context) {
	var req basic_llm_completion.Request
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get completion"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
