package chatgpt

import (
	typings "freechatgpt/internal/typings"
)

func ConvertAPIRequest(api_request typings.APIRequest) typings.ChatGPTRequest {
	chatgpt_request := typings.NewChatGPTRequest(api_request.Model)
	for _, api_message := range api_request.Messages {
		if api_request.Model != "gpt-4" && api_request.Model != "gpt-4-browsing" && api_message.Role == "system" {
			api_message.Role = "critic"
		}
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
	}
	return chatgpt_request
}
