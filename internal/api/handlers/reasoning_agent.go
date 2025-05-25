package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/reasoning_agent"
)

type ReasoningAgentHandler struct {
	service *reasoning_agent.Service
}

func NewReasoningAgentHandler(service *reasoning_agent.Service) *ReasoningAgentHandler {
	return &ReasoningAgentHandler{
		service: service,
	}
}

func (h *ReasoningAgentHandler) HandleReasoningAgentExecution(c *gin.Context) {
	var req reasoning_agent.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Message == "" && req.AgentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either message or agent_id must be provided"})
		return
	}

	resp, err := h.service.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute agent"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
