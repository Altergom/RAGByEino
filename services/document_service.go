package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type DocumentService struct {
	embeddingService *EmbeddingService
}

func NewDocumentService(embeddingService *EmbeddingService) *DocumentService {
	return &DocumentService{embeddingService: embeddingService}
}

func (ds *DocumentService) ProcessDocument(title, content1 string) ([]Document, error) {
	// 文档分割
	chunks := ds.splitText(content1, 500)
	var documents []Document
	for i, chunk := range chunks {
		// 生成唯一ID
		id := ds.generateID(title, i)

		// 获取向量
		vector, err := ds.embeddingService.GetEmbedding(chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to get embedding for chunk %d: %v", i, err)
		}

		documents = append(documents, Document{
			ID:      id,
			Vector:  vector,
			Content: chunk,
		})
	}

	return documents, nil
}

func (ds *DocumentService) generateID(title string, index int) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d_%d", title, index, time.Now().Unix())))
	return hex.EncodeToString(hash[:][:16]) // 取前16位
}

func (ds *DocumentService) splitText(text string, chunkSize int) []string {
	var chunks []string

	// 按句子分割
	sentences := strings.Split(text, "。")
	var currentChunk strings.Builder
	for _, sentence := range sentences {
		if currentChunk.Len()+len(sentence) > chunkSize && currentChunk.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}
		currentChunk.WriteString(sentence)
		currentChunk.WriteString("。")
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}
