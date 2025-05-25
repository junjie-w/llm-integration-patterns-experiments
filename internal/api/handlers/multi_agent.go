package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/multi_agent"
)

type MultiAgentHandler struct {
	service *multi_agent.Service
}

func NewMultiAgentHandler(service *multi_agent.Service) *MultiAgentHandler {
	return &MultiAgentHandler{
		service: service,
	}
}

func (h *MultiAgentHandler) HandleMultiAgentProcess(c *gin.Context) {
	var req multi_agent.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	if req.Message == "" && req.ConversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either message or conversation_id must be provided"})
		return
	}
	
	resp, err := h.service.Process(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process multi-agent request"})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}
