package chatmodels

import (
	"testing"

	"github.com/tmc/langchaingo/llms"
)

func Test_extractModelFromOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []llms.CallOption
		expected string
	}{
		{
			name:     "no options provided",
			options:  []llms.CallOption{},
			expected: "",
		},
		{
			name: "WithModel option provided",
			options: []llms.CallOption{
				llms.WithModel("gemini-3-pro-preview"),
			},
			expected: "gemini-3-pro-preview",
		},
		{
			name: "multiple options including WithModel",
			options: []llms.CallOption{
				llms.WithTemperature(0.7),
				llms.WithModel("gemini-2.5-flash"),
				llms.WithMaxTokens(1000),
			},
			expected: "gemini-2.5-flash",
		},
		{
			name: "other options without WithModel",
			options: []llms.CallOption{
				llms.WithTemperature(0.5),
				llms.WithMaxTokens(500),
			},
			expected: "",
		},
		{
			name: "imagen model for image generation",
			options: []llms.CallOption{
				llms.WithModel("imagen-4.0-generate-001"),
			},
			expected: "imagen-4.0-generate-001",
		},
		{
			name: "gemini flash image preview model",
			options: []llms.CallOption{
				llms.WithModel("gemini-2.5-flash-image-preview"),
			},
			expected: "gemini-2.5-flash-image-preview",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractModelFromOptions(tt.options...)
			if result != tt.expected {
				t.Errorf("extractModelFromOptions() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
