# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an Alexa skill backend that uses **AWS Bedrock** and **Cloudflare Workers AI** for all AI inference. The architecture uses AWS Lambda functions with SQS queues to handle Alexa's 8-second timeout constraint.

## Key Architecture Components

- **Alexa Lambda Handler** (`cmd/alexa/main.go`): Receives Alexa requests, queues them, and polls for responses
- **SQS Request Processor** (`cmd/sqs/main.go`): Processes queued requests using AI providers
- **Queue-based Architecture**: Uses SQS for asynchronous processing to handle Alexa's timeout
- **Three-Provider Design**: Models go through three clients in `internal/dom/chatmodels/`
  - `BedrockApiClient` — Converse API for Claude and Nova models
  - `MantleApiClient` — OpenAI-compatible Responses API (bedrock-mantle endpoint) for Grok and GPT models; holds one SigV4-signed client per region (Grok: `us-west-2`, GPT-5.5: `us-east-1`)
  - `CloudflareApiClient` — OpenAI-compatible Chat Completions API for Llama, Gemma, and Kimi; direct HTTP for Flux image generation; only initialised when `CLOUDFLARE_ACCOUNT_ID` and `CLOUDFLARE_API_KEY` env vars are set

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
# Deploy to AWS — Bedrock uses the Lambda IAM role; Cloudflare params are optional
sam deploy --stack-name alexa-chatgpt \
  --s3-bucket $S3_BUCKET_NAME \
  --parameter-overrides Runtime=provided.al2023 Handler=bootstrap Architecture=arm64 \
    CloudFlareAccountId=$CLOUDFLARE_ACCOUNT_ID CloudFlareAPIKey=$CLOUDFLARE_API_KEY \
  --capabilities CAPABILITY_IAM

# Delete stack
sam delete --stack-name alexa-chatgpt
```

## Code Organization

- **cmd/**: Entry points for Lambda functions
  - `alexa/`: Main Alexa skill handler
  - `sqs/`: Background processor for AI requests
  - `example/`: Local smoke-test binary
- **internal/api/**: Core Alexa request/response handling, game logic
- **internal/dom/chatmodels/**: Bedrock client implementations and model registry
- **internal/pkg/queue/**: SQS queue utilities
- **internal/otel/**: OpenTelemetry instrumentation setup

## Key Implementation Details

### Model Routing

Models are registered in `internal/dom/chatmodels/models.go` with a `Provider` field:

- `ProviderBedrock` → `BedrockApiClient` (Converse API for chat)
- `ProviderBedrockMantle` → `MantleApiClient` (OpenAI Responses API via bedrock-mantle endpoint)
- `ProviderCloudflare` → `CloudflareApiClient` (OpenAI Chat Completions API for chat; direct HTTP for images)

`prompts.go` dispatches to the correct client based on the model's provider. `RegisterAvailableClients(cloudflareAvailable bool)` is called automatically by `NewClient()` — it gates Cloudflare models on whether `CloudflareAPI` is non-nil in the `Resources` struct.

### Adding New AI Models
1. Add model constant in `internal/dom/chatmodels/models.go`
2. Add a `ModelConfig` entry to `allModelConfigs` with the correct `Provider`, `ProviderModelID`, and (for `ProviderBedrockMantle`) `MantleRegion`
3. `NewMantleApiClient` automatically picks up any new `MantleRegion` values; `CloudflareApiClient` needs no changes for new Cloudflare models
4. If the model uses a genuinely new provider, add a client interface in `api.go`, implement it, add it to `Resources` in `service.go`, initialise it in `internal/pkg/init/resources.go`, and add a dispatch case in `prompts.go`

### Alexa Intent Processing Flow
1. Intent received in `internal/api/handler.go:Invoke()`
2. Request queued to SQS if AI processing needed
3. Response polled with timeout handling
4. Fallback to "response will be available shortly" if timeout

### Environment Variables Required
- `S3_BUCKET_NAME`: S3 bucket for SAM deployments
- `REQUESTS_QUEUE_URI`: SQS queue for requests (auto-configured by SAM)
- `RESPONSES_QUEUE_URI`: SQS queue for responses (auto-configured by SAM)
- `CLOUDFLARE_ACCOUNT_ID`: *(optional)* Cloudflare account ID — enables llama, gemma, kimi, flux
- `CLOUDFLARE_API_KEY`: *(optional)* Cloudflare API key — required alongside `CLOUDFLARE_ACCOUNT_ID`

Bedrock auth uses the Lambda IAM execution role (no keys needed). Cloudflare models are silently excluded from the registry if the env vars are absent.

## Testing Approach

- Unit tests exist for core components (handlers, games, utilities)
- Mock interfaces (`mockBedrockAPI`, `mockMantleAPI`, `mockCloudflareAPI`) are defined in `api.go` for testing
- Test events are in `events/` directory for local Lambda testing
- Use `-race` flag to detect race conditions in concurrent code

## Common Issues & Solutions

- **Alexa Timeout**: Responses taking >7 seconds trigger async handling via "LastResponseIntent"
- **Model Switching**: Use "Model <alias>" intent to switch between AI providers
- **Image Generation**: Uses Cloudflare Flux Schnell; output is base64 JPEG decoded from the REST response, resized to two resolutions, and stored in S3
- **Bedrock Model Access**: Enable each model in the AWS Bedrock console under **Model access** before deploying
- **Cloudflare Models Absent**: If `CLOUDFLARE_ACCOUNT_ID`/`CLOUDFLARE_API_KEY` are not set at Lambda startup, all Cloudflare models (llama, gemma, kimi, flux) are excluded — check Lambda env vars if they don't appear in "model available"
