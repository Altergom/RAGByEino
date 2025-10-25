package handlers

import (
	"github.com/Altergom/tryEino/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatHandler struct {
	RAGService *services.RAGService
}

type chatReq struct {
	Question string `json:"question" binding:"required"`
}

func NewChatHandler(ragService *services.RAGService) *ChatHandler {
	return &ChatHandler{RAGService: ragService}
}

func (ch *ChatHandler) AskQuestion(c *gin.Context) {
	var req chatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 获取答案
	answer, err := ch.RAGService.AskQuestion(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate answer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"question": req.Question,
		"answer":   answer,
	})
}
