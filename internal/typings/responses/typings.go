package responses

import "encoding/json"

type MetaDataListItem struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type CitationFormat struct {
	Name string `json:"name"`
}

type CiteMetadata struct {
	CitationFormat CitationFormat     `json:"citation_format"`
	MetaDataList   []MetaDataListItem `json:"metadata_list"`
}

type Message struct {
	ID         string      `json:"id"`
	Author     Author      `json:"author"`
	CreateTime float64     `json:"create_time"`
	UpdateTime interface{} `json:"update_time"`
	Content    Content     `json:"content"`
	EndTurn    interface{} `json:"end_turn"`
	Weight     float64     `json:"weight"`
	Metadata   Metadata    `json:"metadata"`
	Recipient  string      `json:"recipient"`
	Status     string      `json:"status"`
}

type Author struct {
	Role     string                 `json:"role"`
	Name     interface{}            `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
	Language    string   `json:"language"`
	Text        string   `json:"text"`
}

type CitationMetadata struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Text  string `json:"text"`
}

type Citation struct {
	StartIx  int              `json:"start_ix"`
	EndIx    int              `json:"end_ix"`
	Metadata CitationMetadata `json:"metadata"`
}

type Metadata struct {
	Timestamp     string       `json:"timestamp_"`
	MessageType   interface{}  `json:"message_type"`
	FinishDetails interface{}  `json:"finish_details"`
	CiteMetadata  CiteMetadata `json:"_cite_metadata"`
	Citations     []Citation   `json:"citations"`
}

type Data struct {
	Message        Message     `json:"message"`
	ConversationID string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}
type ChatCompletionChunk struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
}

func (chunk *ChatCompletionChunk) String() string {
	resp, _ := json.Marshal(chunk)
	return string(resp)
}

type Choices struct {
	Delta        Delta       `json:"delta"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func NewChatCompletionChunk(model, text string) ChatCompletionChunk {
	return ChatCompletionChunk{
		ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
		Object:  "chat.completion.chunk",
		Created: 0,
		Model:   model,
		Choices: []Choices{
			{
				Index: 0,
				Delta: Delta{
					Content: text,
					Role:    "assistant",
				},
				FinishReason: nil,
			},
		},
	}
}

func StopChunk(model string) ChatCompletionChunk {
	return ChatCompletionChunk{
		ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
		Object:  "chat.completion.chunk",
		Created: 0,
		Model:   model,
		Choices: []Choices{
			{
				Index:        0,
				FinishReason: "stop",
			},
		},
	}
}

type ChatCompletion struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}
type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Choice struct {
	Index        int         `json:"index"`
	Message      Msg         `json:"message"`
	FinishReason interface{} `json:"finish_reason"`
}
type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewChatCompletion(full_test string) ChatCompletion {
	return ChatCompletion{
		ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
		Object:  "chat.completion",
		Created: int64(0),
		Model:   "gpt-3.5-turbo-0301",
		Usage: usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		Choices: []Choice{
			{
				Message: Msg{
					Content: full_test,
					Role:    "assistant",
				},
				Index: 0,
			},
		},
	}
}
