package chatmodels

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	localOtel "github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type BedrockApiClient struct {
	Client *bedrockruntime.Client
}

func NewBedrockApiClient() *BedrockApiClient {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(localOtel.DefaultTransportFormatter),
		),
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithHTTPClient(httpClient),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config for Bedrock: %v", err))
	}
	return &BedrockApiClient{
		Client: bedrockruntime.NewFromConfig(cfg),
	}
}

// GenerateContent calls the Bedrock Converse API.
func (api *BedrockApiClient) GenerateContent(
	ctx context.Context,
	messages []Message,
	opts GenerateOptions,
) (*GenerateResponse, error) {
	var bedrockMessages []bedrocktypes.Message
	var systemPrompts []bedrocktypes.SystemContentBlock

	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			systemPrompts = append(systemPrompts, &bedrocktypes.SystemContentBlockMemberText{
				Value: msg.Content,
			})
		case RoleUser:
			bedrockMessages = append(bedrockMessages, bedrocktypes.Message{
				Role: bedrocktypes.ConversationRoleUser,
				Content: []bedrocktypes.ContentBlock{
					&bedrocktypes.ContentBlockMemberText{Value: msg.Content},
				},
			})
		case RoleAssistant:
			bedrockMessages = append(bedrockMessages, bedrocktypes.Message{
				Role: bedrocktypes.ConversationRoleAssistant,
				Content: []bedrocktypes.ContentBlock{
					&bedrocktypes.ContentBlockMemberText{Value: msg.Content},
				},
			})
		}
	}

	input := &bedrockruntime.ConverseInput{
		ModelId:  aws.String(opts.Model),
		Messages: bedrockMessages,
	}

	if len(systemPrompts) > 0 {
		input.System = systemPrompts
	}

	if opts.Temperature > 0 || opts.MaxTokens > 0 {
		inferenceConfig := &bedrocktypes.InferenceConfiguration{}
		if opts.Temperature > 0 {
			temp := float32(opts.Temperature)
			inferenceConfig.Temperature = &temp
		}
		if opts.MaxTokens > 0 {
			maxTokens := int32(opts.MaxTokens)
			inferenceConfig.MaxTokens = &maxTokens
		}
		input.InferenceConfig = inferenceConfig
	}

	resp, err := api.Client.Converse(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("bedrock converse error: %w", err)
	}

	var responseText strings.Builder
	if outputMsg, ok := resp.Output.(*bedrocktypes.ConverseOutputMemberMessage); ok {
		for _, block := range outputMsg.Value.Content {
			if textBlock, ok := block.(*bedrocktypes.ContentBlockMemberText); ok {
				responseText.WriteString(textBlock.Value)
			}
		}
	}

	return &GenerateResponse{Content: responseText.String()}, nil
}

// GenerateImage calls the Bedrock InvokeModel API.
// Supports Nova Canvas and Titan Image Generator.
func (api *BedrockApiClient) GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error) {
	body, err := json.Marshal(map[string]any{
		"taskType": "TEXT_IMAGE",
		"textToImageParams": map[string]string{
			"text": prompt,
		},
		"imageGenerationConfig": map[string]any{
			"numberOfImages": 1,
			"width":          1024,
			"height":         1024,
			"quality":        "standard",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("bedrock image: failed to marshal request: %w", err)
	}

	resp, err := api.Client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(model),
		Body:        body,
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
	})
	if err != nil {
		return nil, fmt.Errorf("bedrock image invoke error: %w", err)
	}

	var result struct {
		Images []string `json:"images"`
		Error  string   `json:"error"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("bedrock image: failed to unmarshal response: %w", err)
	}
	if result.Error != "" {
		return nil, fmt.Errorf("bedrock image error: %s", result.Error)
	}
	if len(result.Images) == 0 {
		return nil, fmt.Errorf("bedrock image: no images in response")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(result.Images[0])
	if err != nil {
		return nil, fmt.Errorf("bedrock image: failed to decode base64 image: %w", err)
	}
	return imgBytes, nil
}
