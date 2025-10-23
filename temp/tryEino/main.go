package main

import (
	"github.com/Altergom/tryEino/config"
	"github.com/Altergom/tryEino/services"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// godotenv库加载env环境配置
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.InitConfig()

	milvusServier, err := services.NewMilvusService(cfg)
	if err != nil {
		log.Fatal("Error creating milvus service")
	}
	defer milvusServier.Close()

}
