package chatgpt

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	Service
	mock.Mock
}

func (client *MockClient) AutoComplete(ctx context.Context, prompt string) (string, error) {
	args := client.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}
