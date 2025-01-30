package chatmodels

type LastResponse struct {
	Prompt         string   `json:"prompt"`
	Response       string   `json:"response"`
	TimeDiff       string   `json:"time_diff"`
	Model          string   `json:"model"`
	ImagesResponse []string `json:"images_responses"`
	Error          string   `json:"error_message"`
	UserID         string   `json:"user_id"`
}

type Request struct {
	Prompt         string      `json:"prompt"`
	TargetLanguage string      `json:"target_language,omitempty"`
	SourceLanguage string      `json:"source_language,omitempty"`
	Model          ChatModel   `json:"model"`
	ImageModel     *ImageModel `json:"image_model"`
	UserID         string      `json:"user_id"`
}
