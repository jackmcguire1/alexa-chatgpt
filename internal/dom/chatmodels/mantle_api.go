package chatmodels

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
	localOtel "github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// MantleApiClient calls the AWS Bedrock Mantle endpoint using the OpenAI-compatible
// Responses API. Requests are signed with AWS SigV4.
type MantleApiClient struct {
	client openai.Client
}

// sigV4Transport signs each outbound HTTP request with AWS SigV4 before forwarding.
type sigV4Transport struct {
	inner  http.RoundTripper
	signer *v4.Signer
	creds  aws.CredentialsProvider
	region string
}

func (t *sigV4Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	creds, err := t.creds.Retrieve(req.Context())
	if err != nil {
		return nil, fmt.Errorf("mantle: retrieve AWS credentials: %w", err)
	}

	var bodyBytes []byte
	var payloadHash string

	if req.Body != nil && req.Body != http.NoBody {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("mantle: read request body: %w", err)
		}
		req.Body.Close()
		h := sha256.Sum256(bodyBytes)
		payloadHash = hex.EncodeToString(h[:])
	} else {
		// SHA256 of empty string
		payloadHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	}

	reqToSign := req.Clone(req.Context())
	if len(bodyBytes) > 0 {
		reqToSign.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		reqToSign.ContentLength = int64(len(bodyBytes))
	}

	if err = t.signer.SignHTTP(reqToSign.Context(), creds, reqToSign, payloadHash, "bedrock", t.region, time.Now()); err != nil {
		return nil, fmt.Errorf("mantle: sign request: %w", err)
	}

	return t.inner.RoundTrip(reqToSign)
}

func NewMantleApiClient() *MantleApiClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("mantle: failed to load AWS config: %v", err))
	}

	signingTransport := &sigV4Transport{
		inner:  http.DefaultTransport,
		signer: v4.NewSigner(),
		creds:  cfg.Credentials,
		region: cfg.Region,
	}

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			signingTransport,
			otelhttp.WithSpanNameFormatter(localOtel.DefaultTransportFormatter),
		),
	}

	baseURL := fmt.Sprintf("https://bedrock-mantle.%s.api.aws/openai/v1", cfg.Region)

	cl := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithHTTPClient(httpClient),
		option.WithAPIKey("x"), // auth is handled by SigV4; SDK requires a non-empty value
	)

	return &MantleApiClient{client: cl}
}

func (api *MantleApiClient) GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error) {
	params := responses.ResponseNewParams{
		Model: opts.Model,
	}

	var userContent string
	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			params.Instructions = openai.String(msg.Content)
		case RoleUser:
			userContent = msg.Content
		}
	}

	params.Input = responses.ResponseNewParamsInputUnion{
		OfString: openai.String(userContent),
	}

	if opts.Temperature > 0 {
		params.Temperature = openai.Float(opts.Temperature)
	}
	if opts.MaxTokens > 0 {
		params.MaxOutputTokens = openai.Int(int64(opts.MaxTokens))
	}

	resp, err := api.client.Responses.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("mantle: responses API error: %w", err)
	}

	return &GenerateResponse{Content: resp.OutputText()}, nil
}
