package services

import (
	"context"
	"fmt"
	"github.com/Altergom/tryEino/config"
	"github.com/Altergom/tryEino/prompt"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"strings"
)

type RAGService struct {
	embeddingService *EmbeddingService
	milvusService    *MilvusService
	chatModel        *ark.ChatModel
	cfg              *config.Config
}

var RS *RAGService

func NewRAGService(embeddingService *EmbeddingService, milvusService *MilvusService, cfg *config.Config) (*RAGService, error) {
	ctx := context.Background()
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		Model:  cfg.ChatModelVolcano,
		APIKey: cfg.VolcanoAPIKey,
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
	searchResults, _ := s.milvusService.Search(questionVector, s.cfg.TopK)

	// 构建上下文
	context1 := s.buildContext(searchResults)

	// 生成回答
	answer, err := s.generateAnswer(ctx, question, context1, "知识库助手")
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

func (s *RAGService) generateAnswer(ctx context.Context, question, context, role string) (string, error) {
	variables := prompt.SetTemplate(role, fmt.Sprintf("上下文：\n%s\n\n问题：%s\n\n你的回答会用到上下文但不要提及上下文", context, question), []*schema.Message{})
	messages, err := prompt.Template.Format(ctx, variables)
	if err != nil {
		return "", err
	}

	res, err := s.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate answer: %v", err)
	}

	return res.Content, nil
}
