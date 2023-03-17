package queue

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockQueue struct {
	mock.Mock
}

func (q *MockQueue) PushMessage(ctx context.Context, i any) error {
	args := q.Called(ctx, i)
	return args.Error(0)
}

func (q *MockQueue) PullMessage(ctx context.Context, wait int) (res []byte, err error) {
	args := q.Called(ctx, wait)

	if args.Get(0) != nil {
		res = args.Get(0).([]byte)
	}

	return res, args.Error(1)
}
