package main

import (
	"fmt"
	"github.com/Altergom/tryEino/config"
	"github.com/Altergom/tryEino/handlers"
	"github.com/Altergom/tryEino/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// godotenv库加载env环境配置
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.InitConfig()

	// 初始化milvus服务
	milvusService, err := services.NewMilvusService(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize Milvus service: %v\n", err)
	}
	defer milvusService.Close()

	// 初始化embedding服务
	embeddingService, err := services.NewEmbeddingService(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize embedding service: %v\n", err)
	}

	// 初始化文档服务
	documentService := services.NewDocumentService(embeddingService)

	// 初始化RAG服务
	RAGService, err := services.NewRAGService(embeddingService, milvusService, cfg)
	if err != nil {
		fmt.Printf("Failed to initialize RAG service: %v\n", err)
	}

	// 初始化路由
	documentHandler := handlers.NewDocumentHandler(documentService, milvusService)
	chatHandler := handlers.NewChatHandler(RAGService)

	// 设置路由
	router := gin.Default()
	router.POST("/api/documents", documentHandler.UploadDocument)
	router.POST("/api/chat", chatHandler.AskQuestion)

	go func() {
		if err := router.Run(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}
