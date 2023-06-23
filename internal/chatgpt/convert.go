package chatgpt

import (
	typings "freechatgpt/internal/typings"
	"math/rand"
	"strconv"
)

func generate_random_hex(length int) string {
	const charset = "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func randint(min int, max int) int {
	return rand.Intn(max-min) + min
}

func ConvertAPIRequest(api_request typings.APIRequest) typings.ChatGPTRequest {
	chatgpt_request := typings.NewChatGPTRequest(api_request.Model)
	for _, api_message := range api_request.Messages {
		if api_request.Model != "gpt-4" && api_request.Model != "gpt-4-browsing" && api_message.Role == "system" {
			api_message.Role = "critic"
		}
		chatgpt_request.AddMessage(api_message.Role, api_message.Content)
		chatgpt_request.ArkoseToken = generate_random_hex(17) + "|r=ap-southeast-1|meta=3|meta_width=300|metabgclr=transparent|metaiconclr=%23555555|guitextcolor=%23000000|pk=35536E1E-65B4-4D96-9D97-6ADB7EFF8147|at=40|sup=1|rid=" + strconv.Itoa(randint(1, 99)) + "|ag=101|cdn_url=https%3A%2F%2Ftcr9i.chat.openai.com%2Fcdn%2Ffc|lurl=https%3A%2F%2Faudio-ap-southeast-1.arkoselabs.com|surl=https%3A%2F%2Ftcr9i.chat.openai.com|smurl=https%3A%2F%2Ftcr9i.chat.openai.com%2Fcdn%2Ffc%2Fassets%2Fstyle-manager"
	}
	return chatgpt_request
}
