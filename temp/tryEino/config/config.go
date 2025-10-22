package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Milvus配置
	MilvusHost string
	MilvusPort int

	// 应用配置
	Port     string
	LogLevel string

	// 向量配置
	VectorDim      int
	CollectionName string
}

func InitConfig() *Config {
	return &Config{
		MilvusHost:     getEnv("MILVUS_HOST", "localhost"),
		MilvusPort:     getEnvAsInt("MILVUS_PORT", 19530),
		Port:           getEnv("PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		VectorDim:      getEnvAsInt("VECTOR_DIM", 768),
		CollectionName: getEnv("COLLECTION_NAME", "knowledge_base"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}
