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
