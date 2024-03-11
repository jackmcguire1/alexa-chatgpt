package chatmodels

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	Service
	mock.Mock
}

func (client *MockClient) AutoComplete(ctx context.Context, prompt string, model ChatModel) (string, error) {
	args := client.Called(ctx, prompt, model)
	return args.String(0), args.Error(1)
}

func (client *MockClient) GenerateImage(ctx context.Context, prompt string, model ChatModel) (res []byte, err error) {
	args := client.Called(ctx, prompt, model)
	if args.Get(0) != nil {
		res = args.Get(0).([]byte)
	}
	return res, args.Error(1)
}
