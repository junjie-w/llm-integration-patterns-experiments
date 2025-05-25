package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/function_calling"
)

type FunctionCallingHandler struct {
	service *function_calling.Service
}

func NewFunctionCallingHandler(service *function_calling.Service) *FunctionCallingHandler {
	return &FunctionCallingHandler{
		service: service,
	}
}

func (h *FunctionCallingHandler) HandleFunctionCallingCompletion(c *gin.Context) {
	var req function_calling.Request
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get function calling completion"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
