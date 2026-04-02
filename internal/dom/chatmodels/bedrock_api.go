package chatmodels

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/tmc/langchaingo/llms"
)

type BedrockApiClient struct {
	Client *bedrockruntime.Client
}

func NewBedrockApiClient() *BedrockApiClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config for Bedrock: %v", err))
	}
	return &BedrockApiClient{
		Client: bedrockruntime.NewFromConfig(cfg),
	}
}

// GetModel returns itself as it implements llms.Model
func (api *BedrockApiClient) GetModel(options ...llms.CallOption) llms.Model {
	return api
}

// Call implements llms.Model interface
func (api *BedrockApiClient) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	r, err := api.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}, options...)
	if err != nil {
		return "", err
	}
	if len(r.Choices) == 0 {
		return "", fmt.Errorf("no response from Bedrock")
	}
	return r.Choices[0].Content, nil
}

// GenerateContent implements llms.Model and LlmContentGenerator using the Bedrock Converse API
func (api *BedrockApiClient) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	opts := &llms.CallOptions{}
	for _, opt := range options {
		opt(opts)
	}

	var bedrockMessages []bedrocktypes.Message
	var systemPrompts []bedrocktypes.SystemContentBlock

	for _, msg := range messages {
		switch msg.Role {
		case llms.ChatMessageTypeSystem:
			for _, part := range msg.Parts {
				if textPart, ok := part.(llms.TextContent); ok {
					systemPrompts = append(systemPrompts, &bedrocktypes.SystemContentBlockMemberText{
						Value: textPart.Text,
					})
				}
			}
		case llms.ChatMessageTypeHuman:
			var content []bedrocktypes.ContentBlock
			for _, part := range msg.Parts {
				if textPart, ok := part.(llms.TextContent); ok {
					content = append(content, &bedrocktypes.ContentBlockMemberText{
						Value: textPart.Text,
					})
				}
			}
			bedrockMessages = append(bedrockMessages, bedrocktypes.Message{
				Role:    bedrocktypes.ConversationRoleUser,
				Content: content,
			})
		case llms.ChatMessageTypeAI:
			var content []bedrocktypes.ContentBlock
			for _, part := range msg.Parts {
				if textPart, ok := part.(llms.TextContent); ok {
					content = append(content, &bedrocktypes.ContentBlockMemberText{
						Value: textPart.Text,
					})
				}
			}
			bedrockMessages = append(bedrockMessages, bedrocktypes.Message{
				Role:    bedrocktypes.ConversationRoleAssistant,
				Content: content,
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

	var responseText string
	if outputMsg, ok := resp.Output.(*bedrocktypes.ConverseOutputMemberMessage); ok {
		for _, block := range outputMsg.Value.Content {
			if textBlock, ok := block.(*bedrocktypes.ContentBlockMemberText); ok {
				responseText += textBlock.Value
			}
		}
	}

	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{
			{Content: responseText},
		},
	}, nil
}

// GenerateImage implements BedrockAPI using the InvokeModel API.
// Supports Nova Canvas and Titan Image Generator (both share the same request format).
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