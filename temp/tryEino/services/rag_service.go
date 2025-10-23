package services

import (
	"context"
	"fmt"
	"github.com/Altergom/tryEino/config"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"strings"
)

type RAGService struct {
	embeddingService *EmbeddingService
	milvusService    *MilvusService
	chatModel        *openai.ChatModel
	cfg              *config.Config
}

func NewRAGService(embeddingService *EmbeddingService, milvusService *MilvusService, cfg *config.Config) (*RAGService, error) {
	ctx := context.Background()
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:  "gpt-3.5-turbo",
		APIKey: cfg.OpenAIAPIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chat model: %v", err)
	}

	return &RAGService{
		embeddingService: embeddingService,
		milvusService:    milvusService,
		chatModel:        chatModel,
		cfg:              cfg,
	}, nil
}

func (s *RAGService) AskQuestion(question string) (string, error) {
	ctx := context.Background()

	// 问题向量化
	questionVector, err := s.embeddingService.GetEmbedding(question)
	if err != nil {
		return "", fmt.Errorf("failed to get question embedding: %v", err)
	}

	// 在milvus搜索相似文档
	searchResults, err := s.milvusService.Search(questionVector, s.cfg.TopK)
	if err != nil {
		return "", fmt.Errorf("failed to search documents: %v", err)
	}

	// 构建上下文
	context1 := s.buildContext(searchResults)

	// 生成回答
	answer, err := s.generateAnswer(ctx, question, context1)
	if err != nil {
		return "", fmt.Errorf("failed to generate answer: %v", err)
	}

	return answer, nil
}

func (s *RAGService) buildContext(searchResults []SearchResult) string {
	var contextParts []string
	for _, result := range searchResults {
		contextParts = append(contextParts, result.Content)
	}

	return strings.Join(contextParts, "\n\n")
}

func (s *RAGService) generateAnswer(ctx context.Context, question string, context1 string) (string, error) {
	messages := []*schema.Message{
		schema.SystemMessage("你是一个知识库助手，请基于提供的上下文信息回答问题"),
		schema.UserMessage(fmt.Sprintf("上下文：\n%s\n\n问题：%s", context1, question)),
	}

	res, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate answer: %v", err)
	}

	return res.Content, nil
}
