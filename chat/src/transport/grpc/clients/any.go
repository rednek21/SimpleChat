package clients

import (
	"context"
)

type AnyClient struct {
}

func NewAnyClient() *AnyClient {
	return &AnyClient{}
}

func (l *AnyClient) SendMessage(ctx context.Context, msg string) (string, error) {
	return msg, nil
}
