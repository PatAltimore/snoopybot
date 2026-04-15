package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

const (
	tableName    = "state"
	partitionKey = "state"
	rowKey       = "novel"
)

type StateClient struct {
	tableClient *aztables.Client
	svcClient   *aztables.ServiceClient
}

type novelEntity struct {
	PartitionKey string `json:"PartitionKey"`
	RowKey       string `json:"RowKey"`
	Index        int    `json:"index"`
}

func NewStateClient(account, key string) (*StateClient, error) {
	cred, err := aztables.NewSharedKeyCredential(account, key)
	if err != nil {
		return nil, err
	}

	svcURL := fmt.Sprintf("https://%s.table.core.windows.net/", account)
	svcClient, err := aztables.NewServiceClientWithSharedKey(svcURL, cred, nil)
	if err != nil {
		return nil, err
	}

	tableClient := svcClient.NewClient(tableName)

	return &StateClient{
		tableClient: tableClient,
		svcClient:   svcClient,
	}, nil
}

func (s *StateClient) EnsureTable(ctx context.Context) error {
	_, err := s.svcClient.CreateTable(ctx, tableName, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == http.StatusConflict {
			// Table already exists — fine
		} else {
			return err
		}
	}

	_, err = s.tableClient.GetEntity(ctx, partitionKey, rowKey, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == http.StatusNotFound {
			entity := novelEntity{PartitionKey: partitionKey, RowKey: rowKey, Index: 0}
			data, merr := json.Marshal(entity)
			if merr != nil {
				return merr
			}
			_, merr = s.tableClient.AddEntity(ctx, data, nil)
			return merr
		}
		return err
	}
	return nil
}

func (s *StateClient) GetNovelIndex(ctx context.Context) int {
	resp, err := s.tableClient.GetEntity(ctx, partitionKey, rowKey, nil)
	if err != nil {
		return 0
	}
	var entity novelEntity
	if err := json.Unmarshal(resp.Value, &entity); err != nil {
		return 0
	}
	return entity.Index
}

func (s *StateClient) SetNovelIndex(ctx context.Context, index int) error {
	entity := novelEntity{PartitionKey: partitionKey, RowKey: rowKey, Index: index}
	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	_, err = s.tableClient.UpdateEntity(ctx, data, &aztables.UpdateEntityOptions{
		UpdateMode: aztables.UpdateModeMerge,
	})
	return err
}
