package es

import (
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type CanalData struct {
	Data      json.RawMessage  `json:"data"`
	Database  string           `json:"database"`
	Es        int32            `json:"es"`
	ID        int32            `json:"id"`
	IsDdl     bool             `json:"isDdl"`
	MysqlType map[string]any   `json:"mysqlType"`
	Old       []map[string]any `json:"old"`
	PkNames   []string         `json:"pkNames"`
	Sql       string           `json:"sql"`
	SqlType   map[string]any   `json:"sqlType"`
	Table     string           `json:"table"`
	Ts        int64            `json:"ts"`
	Type      string           `json:"type"`
}

func GoodsSaveToES(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var data CanalData
		_ = json.Unmarshal(msgs[i].Body, &data)
		if data.Type == "UPDATE" {
			// TODO

		}
	}

	return consumer.ConsumeRetryLater, nil
}
