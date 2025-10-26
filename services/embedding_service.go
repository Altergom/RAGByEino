package services

import (
	"context"
	"fmt"
	"github.com/Altergom/tryEino/config"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"time"
)

type EmbeddingService struct {
	embedder *ark.Embedder
	cfg      *config.Config
}

var ES *EmbeddingService

func NewEmbeddingService(cfg *config.Config) (*EmbeddingService, error) {
	ctx := context.Background()

	timeout := time.Second * 30
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  cfg.VolcanoAPIKey,
		Model:   cfg.EmbeddingModelVolcano,
		Timeout: &timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Ark embedder: %v", err)
	}

	return &EmbeddingService{
		embedder: embedder,
		cfg:      cfg,
	}, nil
}

func (es *EmbeddingService) GetEmbedding(text string) ([]float32, error) {
	ctx := context.Background()

	embeddings, err := es.embedder.EmbedStrings(ctx, []string{text})
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding: %v", err)
	}

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embedding data received")
	}

	res := make([]float32, len(embeddings[0]))
	for i, embedding := range embeddings[0] {
		res[i] = float32(embedding)
	}

	return res, nil
}

func (es *EmbeddingService) GetEmbeddings(texts []string) ([][]float32, error) {
	ctx := context.Background()

	embeddings, err := es.embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to get embeddings: %v", err)
	}

	// 转换为 float32 切片
	result := make([][]float32, len(embeddings))
	for i, embedding := range embeddings {
		resultVector := make([]float32, len(embedding))
		for j, v := range embedding {
			resultVector[j] = float32(v)
		}
		result[i] = resultVector
	}

	return result, nil
}
