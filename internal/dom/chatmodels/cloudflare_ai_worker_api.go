package chatmodels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

const (
	CF_LLAMA_2_7B_CHAT_INT8_MODEL = "@cf/meta/llama-2-7b-chat-int8"
	CF_LLAMA_3_8B_INSTRUCT_MODEL  = "@cf/meta/llama-3-8b-instruct"
	CF_SQL_MODEL                  = "@cf/defog/sqlcoder-7b-2"
	CF_AWQ_MODEL                  = "@hf/thebloke/llama-2-13b-chat-awq"
	CF_OPEN_CHAT_MODEL            = "@cf/openchat/openchat-3.5-0106"
	CF_STABLE_DIFFUSION           = "@cf/stabilityai/stable-diffusion-xl-base-1.0"
	CF_META_TRANSLATION_MODEL     = "@cf/meta/m2m100-1.2b"
	CF_QWEN_MODEL                 = "@cf/qwen/qwen1.5-1.8b-chat"
)

var CHAT_MODEL_TO_CF_MODEL = map[ChatModel]string{
	CHAT_MODEL_SQL:              CF_SQL_MODEL,
	CHAT_MODEL_AWQ:              CF_AWQ_MODEL,
	CHAT_MODEL_META:             CF_LLAMA_3_8B_INSTRUCT_MODEL,
	CHAT_MODEL_OPEN:             CF_OPEN_CHAT_MODEL,
	CHAT_MODEL_STABLE_DIFFUSION: CF_STABLE_DIFFUSION,
	CHAT_MODEL_TRANSLATIONS:     CF_META_TRANSLATION_MODEL,
	CHAT_MODEL_QWEN:             CF_QWEN_MODEL,
}

type Response struct {
	Result struct {
		Response string `json:"response"`
	} `json:"result,omitempty"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
	Messages []string `json:"messages,omitempty"`
	Success  bool     `json:"success"`
}

type TranslateResponse struct {
	Result struct {
		TranslatedText string `json:"translated_text"`
	} `json:"result"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
	Messages []string `json:"messages,omitempty"`
	Success  bool     `json:"success"`
}

type CloudflareApiClient struct {
	AccountID string
	APIKey    string
}

func NewCloudflareApiClient(accountID, apiKey string) *CloudflareApiClient {
	return &CloudflareApiClient{
		AccountID: accountID,
		APIKey:    apiKey,
	}
}

func (api *CloudflareApiClient) GenerateText(ctx context.Context, prompt string, model string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", api.AccountID, model)

	payload := map[string]string{
		"prompt": prompt,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+api.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result *Response
	err = json.Unmarshal(data, &result)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal cloudflare response body %s", string(data))
		return "", err
	}

	if !result.Success {
		err = fmt.Errorf("didn't get success from result %v http-status: %d", utils.ToJSON(result), resp.StatusCode)
		return "", err
	}

	return result.Result.Response, nil
}

func (api *CloudflareApiClient) GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", api.AccountID, model)

	payload := map[string]string{
		"prompt": prompt,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+api.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type GenerateTranslationRequest struct {
	SourceLanguage string
	TargetLanguage string
	Prompt         string
	Model          string
}

func (api *CloudflareApiClient) GenerateTranslation(ctx context.Context, req *GenerateTranslationRequest) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", api.AccountID, req.Model)

	if req.SourceLanguage == "" {
		req.SourceLanguage = "en"
	}
	payload := map[string]string{
		"text":        req.Prompt,
		"source_lang": req.SourceLanguage,
		"target_lang": req.TargetLanguage,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+api.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	slog.
		With("response", string(data)).
		Info("generated translation response")

	var response *TranslateResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}

	if !response.Success {
		err = fmt.Errorf("didn't get success from result %v http-status: %d", utils.ToJSON(response), resp.StatusCode)
		return "", err
	}

	return response.Result.TranslatedText, nil
}
