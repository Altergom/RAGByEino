package main

import (
	"context"
	"fmt"
	"github.com/Altergom/tryEino/config"
	"github.com/Altergom/tryEino/handlers"
	"github.com/Altergom/tryEino/prompt"
	"github.com/Altergom/tryEino/services"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// godotenv库加载env环境配置
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Cfg = config.InitConfig()

	// 初始化milvus服务
	services.MS, err = services.NewMilvusService(config.Cfg)
	if err != nil {
		fmt.Printf("Failed to initialize Milvus service: %v\n", err)
	}
	defer services.MS.Close()

	// 初始化embedding服务
	services.ES, err = services.NewEmbeddingService(config.Cfg)
	if err != nil {
		fmt.Printf("Failed to initialize embedding service: %v\n", err)
	}

	// 初始化文档服务
	services.DS = services.NewDocumentService(services.ES)

	// 初始化模板
	prompt.Template = prompt.NewTemplate()
	tem1 := prompt.SetTemplate("面试官", "", []*schema.Message{})
	prompt.Template.Format(context.Background(), tem1)

	// 初始化RAG服务
	services.RS, err = services.NewRAGService(services.ES, services.MS, config.Cfg)
	if err != nil {
		fmt.Printf("Failed to initialize RAG service: %v\n", err)
	}

	handlers.Run()
}
