package services

import (
	"context"
	"fmt"

	"github.com/Altergom/tryEino/config"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusService struct {
	client client.Client
	cfg    *config.Config
}

type Document struct {
	ID      string
	Content string
	Vector  []float32
}

type SearchResult struct {
	ID      string
	Content string
	Score   float32
}

var MS *MilvusService

func NewMilvusService(cfg *config.Config) (*MilvusService, error) {
	ctx := context.Background()
	conn, err := client.NewGrpcClient(ctx, fmt.Sprintf("%s:%d", cfg.MilvusHost, cfg.MilvusPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %v", err)
	}

	service := &MilvusService{
		client: conn,
		cfg:    cfg,
	}

	err = service.initCollection()
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (ms *MilvusService) initCollection() error {
	ctx := context.Background()

	// 检查集合是否存在
	exist, err := ms.client.HasCollection(ctx, ms.cfg.CollectionName)
	if err != nil {
		return err
	}
	if !exist {
		// 创建集合
		schema := &entity.Schema{
			CollectionName: ms.cfg.CollectionName,
			Description:    "个人知识库",
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeVarChar,
					PrimaryKey: true,
					TypeParams: map[string]string{
						"max_length": "100",
					},
				},
				{
					Name:     "content",
					DataType: entity.FieldTypeVarChar,
					TypeParams: map[string]string{
						"max_length": "65535",
					},
				},
				{
					Name:     "vector",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						"dim": fmt.Sprintf("%d", ms.cfg.VectorDim),
					},
				},
			},
		}
		err = ms.client.CreateCollection(ctx, schema, entity.DefaultShardNumber)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ms *MilvusService) Close() error {
	return ms.client.Close()
}

func (ms *MilvusService) InsertDocument(doc []Document) error {
	ctx := context.Background()
	ids := make([]string, len(doc))
	contents := make([]string, len(doc))
	vector := make([][]float32, len(doc))

	for i, d := range doc {
		ids[i] = d.ID
		contents[i] = d.Content
		vector[i] = d.Vector
	}

	_, err := ms.client.Insert(ctx, ms.cfg.CollectionName, "",
		entity.NewColumnVarChar("id", ids),
		entity.NewColumnVarChar("content", contents),
		entity.NewColumnFloatVector("vector", ms.cfg.VectorDim, vector),
	)
	if err != nil {
		return err
	}

	return nil
}

func (ms *MilvusService) Search(queryVector []float32, topK int) ([]SearchResult, error) {
	ctx := context.Background()
	// 创建搜索参数
	searchParam, _ := entity.NewIndexIvfFlatSearchParam(10)
	// 开始搜索
	vectors := []entity.Vector{entity.FloatVector(queryVector)}
	searchRes, err := ms.client.Search(ctx, ms.cfg.CollectionName, []string{}, "",
		[]string{"id", "content"}, vectors, "vector", entity.L2, topK, searchParam)

	if err != nil {
		return nil, err
	}

	// 处理搜索结果
	var results []SearchResult
	for _, searchResult := range searchRes {
		// 获取ID列
		ids := searchResult.IDs
		// 获取字段数据
		fields := searchResult.Fields
		// 获取分数
		scores := searchResult.Scores

		// 遍历每个结果
		for i := 0; i < searchResult.ResultCount; i++ {
			// 获取第i个结果
			id, _ := ids.GetAsString(i)
			content, _ := fields.GetColumn("content").GetAsString(i)
			score := scores[i]

			results = append(results, SearchResult{
				ID:      id,
				Content: content,
				Score:   score,
			})
		}
	}

	return results, nil
}
