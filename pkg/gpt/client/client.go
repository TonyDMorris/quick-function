package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"time"

	"net/http"

	"github.com/TonyDMorris/quick-function/pkg/gpt/models"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	OpenAIURL = "https://api.openai.com/v1/chat/completions"
	GPT4Model = "gpt-4"
	GPT3Model = "gpt-3.5-turbo"
)

type ChatClientInterface interface {
	Chat(messages []models.Message) (*models.CompletionResponse, error)
}

type ChatClient struct {
	client *retryablehttp.Client
	apiKey string
}

func NewChatClient(apiKey string) *ChatClient {
	retriableClient := retryablehttp.NewClient()
	retriableClient.RetryMax = 5
	retriableClient.HTTPClient.Timeout = time.Minute * 5
	return &ChatClient{
		client: retriableClient,
		apiKey: apiKey,
	}
}

func (c *ChatClient) Chat(messages []models.Message) (*models.CompletionResponse, error) {
	requestBody := models.CompletionRequest{
		Model:    GPT3Model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest(http.MethodPost, OpenAIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	var completionResponse models.CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResponse); err != nil {
		return nil, err
	}

	return &completionResponse, nil
}
