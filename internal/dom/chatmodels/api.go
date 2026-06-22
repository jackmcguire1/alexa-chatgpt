package chatmodels

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MessageRole identifies the sender of a message.
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// Message is a single turn in a conversation.
type Message struct {
	Role    MessageRole
	Content string
}

// GenerateOptions configures a content generation call.
type GenerateOptions struct {
	Model        string
	Temperature  float64
	MaxTokens    int
	MantleRegion string // only used for ProviderBedrockMantle models
}

// GenerateResponse holds the result of a generation call.
type GenerateResponse struct {
	Content string
}

// BedrockAPI is the single interface for all AI operations via AWS Bedrock.
type BedrockAPI interface {
	GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error)
	GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error)
}

// MantleAPI is the interface for chat operations via the AWS Bedrock Mantle endpoint
// (OpenAI-compatible API for third-party models such as xAI Grok and OpenAI GPT).
type MantleAPI interface {
	GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error)
}

// CloudflareAPI is the interface for chat and image operations via Cloudflare Workers AI.
type CloudflareAPI interface {
	GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error)
	GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error)
}

type mockBedrockAPI struct {
	mock.Mock
}

func (m *mockBedrockAPI) GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error) {
	args := m.Called(ctx, messages, opts)
	if args.Get(0) != nil {
		return args.Get(0).(*GenerateResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockBedrockAPI) GenerateImage(ctx context.Context, prompt string, model string) (res []byte, err error) {
	args := m.Called(ctx, prompt, model)
	if args.Get(0) != nil {
		res = args.Get(0).([]byte)
	}
	return res, args.Error(1)
}

type mockMantleAPI struct {
	mock.Mock
}

func (m *mockMantleAPI) GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error) {
	args := m.Called(ctx, messages, opts)
	if args.Get(0) != nil {
		return args.Get(0).(*GenerateResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockCloudflareAPI struct {
	mock.Mock
}

func (m *mockCloudflareAPI) GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error) {
	args := m.Called(ctx, messages, opts)
	if args.Get(0) != nil {
		return args.Get(0).(*GenerateResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockCloudflareAPI) GenerateImage(ctx context.Context, prompt string, model string) (res []byte, err error) {
	args := m.Called(ctx, prompt, model)
	if args.Get(0) != nil {
		res = args.Get(0).([]byte)
	}
	return res, args.Error(1)
}
