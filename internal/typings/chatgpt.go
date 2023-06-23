package types

import "github.com/google/uuid"

type Chatgpt_message struct {
	ID      uuid.UUID       `json:"id"`
	Author  chatgpt_author  `json:"author"`
	Content chatgpt_content `json:"content"`
}

type chatgpt_content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type chatgpt_author struct {
	Role string `json:"role"`
}

type ChatGPTRequest struct {
	Action                     string            `json:"action"`
	Messages                   []Chatgpt_message `json:"messages"`
	ParentMessageID            string            `json:"parent_message_id,omitempty"`
	ConversationID             string            `json:"conversation_id,omitempty"`
	Model                      string            `json:"model"`
	HistoryAndTrainingDisabled bool              `json:"history_and_training_disabled"`
	ArkoseToken                string            `json:"arkose_token,omitempty"`
}

type ChatGPTContinueRequest struct {
	Action          string `json:"action"`
	ParentMessageID string `json:"parent_message_id,omitempty"`
	ConversationID  string `json:"conversation_id,omitempty"`
	Model           string `json:"model"`
}

func NewChatGPTRequest(model string) ChatGPTRequest {
	return ChatGPTRequest{
		Action:          "next",
		ParentMessageID: uuid.NewString(),
		Model:           model,
	}
}

func (c *ChatGPTRequest) AddMessage(role string, content string) {
	c.Messages = append(c.Messages, Chatgpt_message{
		ID:      uuid.New(),
		Author:  chatgpt_author{Role: role},
		Content: chatgpt_content{ContentType: "text", Parts: []string{content}},
	})
}
