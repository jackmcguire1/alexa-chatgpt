package chatmodels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

const (
	CF_LLAMA_2_7B_CHAT_INT8_MODEL = "@cf/meta/llama-2-7b-chat-int8"
	CF_SQL_MODEL                  = "@cf/defog/sqlcoder-7b-2"
	CF_AWQ_MODEL                  = "@hf/thebloke/llama-2-13b-chat-awq"
	CF_OPEN_CHAT_MODEL            = "@cf/openchat/openchat-3.5-0106"
)

var CHAT_MODEL_TO_CF_MODEL = map[ChatModel]string{
	CHAT_MODEL_SQL:  CF_SQL_MODEL,
	CHAT_MODEL_AWQ:  CF_AWQ_MODEL,
	CHAT_MODEL_META: CF_LLAMA_2_7B_CHAT_INT8_MODEL,
	CHAT_MODEL_OPEN: CF_OPEN_CHAT_MODEL,
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
