package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"net/http"

	"github.com/TonyDMorris/quick-function/pkg/gpt/models"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkoukk/tiktoken-go"
)

const MaxTokens = 2048

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
		Model:    GPT4Model,
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

func (c *ChatClient) NumTokensFromMessages(messages models.CompletionRequest, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	var tokensPerMessage int
	switch model {
	case "gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-4-0613",
		"gpt-4-32k-0613":
		tokensPerMessage = 3
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
	default:
		if strings.Contains(model, "gpt-3.5-turbo") {
			log.Println("warning: gpt-3.5-turbo may update over time. Returning num tokens assuming gpt-3.5-turbo-0613.")
			return c.NumTokensFromMessages(messages, "gpt-3.5-turbo-0613")
		} else if strings.Contains(model, "gpt-4") {
			log.Println("warning: gpt-4 may update over time. Returning num tokens assuming gpt-4-0613.")
			return c.NumTokensFromMessages(messages, "gpt-4-0613")
		} else {
			err = fmt.Errorf("num_tokens_from_messages() is not implemented for model %s. See https://github.com/openai/openai-python/blob/main/chatml.md for information on how messages are converted to tokens.", model)
			log.Println(err)
			return
		}
	}

	for _, message := range messages.Messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))

	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
