package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junjie-w/llm-integration-patterns-experiments/internal/ai/evaluation"
)

type EvaluationHandler struct {
	service *evaluation.Service
}

func NewEvaluationHandler(service *evaluation.Service) *EvaluationHandler {
	return &EvaluationHandler{
		service: service,
	}
}

func (h *EvaluationHandler) HandleEvaluate(c *gin.Context) {
	var req evaluation.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query cannot be empty"})
		return
	}
	
	if len(req.PatternTypes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one pattern type must be specified"})
		return
	}
	
	resp, err := h.service.Evaluate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate patterns"})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

func (h *EvaluationHandler) HandleGetReport(c *gin.Context) {
	reportID := c.Param("id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}
	
	report, err := h.service.GetReport(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	
	c.JSON(http.StatusOK, report)
}
