# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an Alexa skill backend that integrates with multiple AI providers (OpenAI, Google Gemini, Anthropic Claude, Cloudflare AI) to enable natural conversations through Alexa devices. The architecture uses AWS Lambda functions with SQS queues to handle Alexa's 8-second timeout constraint.

## Key Architecture Components

- **Alexa Lambda Handler** (`cmd/alexa/main.go`): Receives Alexa requests, queues them, and polls for responses
- **SQS Request Processor** (`cmd/sqs/main.go`): Processes queued requests using AI providers
- **Queue-based Architecture**: Uses SQS for asynchronous processing to handle Alexa's timeout
- **Multi-Provider Support**: Abstracts different AI providers through a common interface in `internal/dom/chatmodels/`

## Essential Commands

### Build & Test
```bash
# Run all tests with race detection and coverage
go test ./... -race -coverprofile=coverage.out -covermode=atomic

# Run tests for a specific package
go test ./internal/api/... -v

# Build SAM application for deployment
sam build --parameter-overrides Runtime=provided.al2023 Handler=bootstrap Architecture=arm64

# Build locally (for ARM64 Lambda)
GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/alexa/main.go
GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/sqs/main.go
```

### Local Development & Testing
```bash
# Start local Lambda runtime for testing
sam local start-lambda

# Invoke function locally with test event
sam local invoke ChatGPTFunction -e events/alexa_request.json

# View Lambda logs locally
sam logs -n ChatGPTLambda --stack-name alexa-chatgpt --tail
```

### Deployment
```bash
# Deploy to AWS (requires configured environment variables)
sam deploy --stack-name alexa-chatgpt \
  --s3-bucket $S3_BUCKET_NAME \
  --parameter-overrides \
    Runtime=provided.al2023 \
    Handler=bootstrap \
    Architecture=arm64 \
    OpenAIApiKey=$OPENAI_API_KEY \
    GeminiApiKey=$GEMINI_API_KEY \
    AnthropicApiKey=$ANTHROPIC_API_KEY \
    CloudflareAccountId=$CLOUDFLARE_ACCOUNT_ID \
    CloudflareApiKey=$CLOUDFLARE_API_KEY \
  --capabilities CAPABILITY_IAM

# Delete stack
sam delete --stack-name alexa-chatgpt
```

## Code Organization

- **cmd/**: Entry points for Lambda functions
  - `alexa/`: Main Alexa skill handler
  - `sqs/`: Background processor for AI requests
- **internal/api/**: Core Alexa request/response handling, game logic
- **internal/dom/chatmodels/**: AI provider abstractions and implementations
- **internal/pkg/queue/**: SQS queue utilities
- **internal/otel/**: OpenTelemetry instrumentation setup

## Key Implementation Details

### Adding New AI Models
1. Define model constants in `internal/dom/chatmodels/models.go`
2. Implement provider client if needed in `internal/dom/chatmodels/`
3. Add model handling in `internal/api/model.go`
4. Update model selection logic in handler

### Alexa Intent Processing Flow
1. Intent received in `internal/api/handler.go:Invoke()`
2. Request queued to SQS if AI processing needed
3. Response polled with timeout handling
4. Fallback to "response will be available shortly" if timeout

### Environment Variables Required
- `OPENAI_API_KEY`: OpenAI API key
- `ANTHROPIC_API_KEY`: Anthropic API key  
- `CLOUDFLARE_ACCOUNT_ID`: Cloudflare account ID
- `CLOUDFLARE_API_KEY`: Cloudflare API key
- `GEMINI_API_KEY`: Base64 encoded Google service account JSON
- `S3_BUCKET_NAME`: S3 bucket for SAM deployments
- `REQUESTS_QUEUE_URI`: SQS queue for requests (auto-configured by SAM)
- `RESPONSES_QUEUE_URI`: SQS queue for responses (auto-configured by SAM)

## Testing Approach

- Unit tests exist for core components (handlers, games, utilities)
- Mock interfaces are generated for testing (e.g., `mock_service.go`)
- Test events are in `events/` directory for local Lambda testing
- Use `-race` flag to detect race conditions in concurrent code

## Common Issues & Solutions

- **Alexa Timeout**: Responses taking >7 seconds trigger async handling via "LastResponseIntent"
- **Model Switching**: Use "Model <alias>" intent to switch between AI providers
- **Image Generation**: Uses S3 for storing generated images, returns pre-signed URLs