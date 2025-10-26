package handlers

import (
	"github.com/Altergom/tryEino/config"
	"github.com/Altergom/tryEino/services"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	// 初始化路由
	documentHandler := NewDocumentHandler(services.DS, services.MS)
	chatHandler := NewChatHandler(services.RS)

	// 设置路由
	router := gin.Default()
	router.POST("/api/documents", documentHandler.UploadDocument)
	router.POST("/api/chat", chatHandler.AskQuestion)

	if err := router.Run(":" + config.Cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
