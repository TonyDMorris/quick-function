package models

const (
	// Roles
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"

	// Finish Reasons
	FinishReasonStop      = "stop"
	FinishReasonMaxLength = "max_length"
	FinishReasonMaxTokens = "max_tokens"
)

// Message represents a single message in a conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest represents the request body for a ChatGPT API call.
type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// CompletionResponseChoice represents a choice in the completion response.
type CompletionResponseChoice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

// CompletionResponse represents the response body from a ChatGPT API call.
type CompletionResponse struct {
	Choices []CompletionResponseChoice `json:"choices"`
}
