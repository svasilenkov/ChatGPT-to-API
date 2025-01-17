package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"freechatgpt/internal/chatgpt"
	"freechatgpt/internal/tokens"
	typings "freechatgpt/internal/typings"
	"freechatgpt/internal/typings/responses"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func passwordHandler(c *gin.Context) {
	// Get the password from the request (json) and update the password
	type password_struct struct {
		Password string `json:"password"`
	}
	var password password_struct
	err := c.BindJSON(&password)
	if err != nil {
		c.String(400, "password not provided")
		return
	}
	ADMIN_PASSWORD = password.Password
	// Set environment variable
	os.Setenv("ADMIN_PASSWORD", ADMIN_PASSWORD)
	c.String(200, "password updated")
}

func tokensHandler(c *gin.Context) {
	// Get the request_tokens from the request (json) and update the request_tokens
	var request_tokens []string
	err := c.BindJSON(&request_tokens)
	if err != nil {
		c.String(400, "tokens not provided")
		return
	}
	ACCESS_TOKENS = tokens.NewAccessToken(request_tokens)
	c.String(200, "tokens updated")
}
func optionsHandler(c *gin.Context) {
	// Set headers for CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
func nightmare(c *gin.Context) {
	var original_request typings.APIRequest
	err := c.BindJSON(&original_request)
	if err != nil {
		c.JSON(400, gin.H{"error": gin.H{
			"message": "Request must be proper JSON",
			"type":    "invalid_request_error",
			"param":   nil,
			"code":    err.Error(),
		}})
	}
	// Convert the chat request to a ChatGPT request
	translated_request := chatgpt.ConvertAPIRequest(original_request)
	// c.JSON(200, chatgpt_request)

	// authHeader := c.GetHeader("Authorization")
	token := ACCESS_TOKENS.GetToken()
	// if authHeader != "" {
	// 	customAccessToken := strings.Replace(authHeader, "Bearer ", "", 1)
	// 	if customAccessToken != "" {
	// 		token = customAccessToken
	// 		println("customAccessToken set:" + customAccessToken)
	// 	}
	// }

	response, err := chatgpt.SendRequest(translated_request, token)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"error": "error sending request" + err.Error(),
		})
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var error_response map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&error_response)
		if err != nil {
			c.JSON(500, gin.H{"error": gin.H{
				"message": "Unknown error: " + err.Error(),
				"type":    "internal_server_error",
				"param":   nil,
				"code":    "500",
			}})
			return
		}
		c.JSON(response.StatusCode, gin.H{"error": gin.H{
			"message": error_response["detail"],
			"type":    response.Status,
			"param":   nil,
			"code":    "error",
		}})
		return
	}
	// Create a bufio.Reader from the response body
	reader := bufio.NewReader(response.Body)

	var fulltext string
	lastMessageId := ""
	var last_browser_metadata responses.Metadata

	// Read the response byte by byte until a newline character is encountered
	if original_request.Stream {
		// Response content type is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content type is application/json
		c.Header("Content-Type", "application/json")
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if len(line) < 6 {
			continue
		}
		// Remove "data: " from the beginning of the line
		line = line[6:]
		if strings.Contains(line, `"status": "finished_successfully"`) {
			//fmt.Println(line + " ")
		}
		var original_response responses.Data

		// Check if line starts with [DONE]
		if strings.HasPrefix(line, "[DONE]") {
			if original_response.Message.EndTurn == false {
				// Try to continue
				translated_request.Action = "continue"
				translated_request.ParentMessageID = lastMessageId
				translated_request.ConversationID = original_response.ConversationID
				translated_request.Messages = nil
				response, err = chatgpt.SendRequest(translated_request, token)
				if err != nil {
					fmt.Println(err.Error())
					c.JSON(500, gin.H{
						"error": "error sending request" + err.Error(),
					})
					return
				}
				defer response.Body.Close()
				if response.StatusCode != 200 {
					// Try read response body as JSON
					var error_response map[string]interface{}
					err = json.NewDecoder(response.Body).Decode(&error_response)
					if err != nil {
						c.JSON(500, gin.H{"error": gin.H{
							"message": "Unknown error: " + err.Error(),
							"type":    "internal_server_error",
							"param":   nil,
							"code":    "500",
						}})
						return
					}
					c.JSON(response.StatusCode, gin.H{"error": gin.H{
						"message": error_response["detail"],
						"type":    response.Status,
						"param":   nil,
						"code":    "error",
					}})
					return
				}
				//bytes, err := ioutil.ReadAll(response.Body)
				//fmt.Println(string(bytes))

				// Create a bufio.Reader from the response body
				reader = bufio.NewReader(response.Body)
				continue
			}
			if !original_request.Stream {
				full_response := responses.NewChatCompletion(fulltext)
				if err != nil {
					return
				}
				c.JSON(200, full_response)
				return
			}
			final_line := responses.StopChunk(original_request.Model)
			c.Writer.WriteString("data: " + final_line.String() + "\n\n")
			//fmt.Println(final_line.String() + "\n")

			c.String(200, "data: [DONE]\n\n")
			return

		} else {
			// Parse the line as JSON
			err = json.Unmarshal([]byte(line), &original_response)
			if err != nil {
				continue
			}
			isNewMesssage := false
			if lastMessageId != original_response.Message.ID {
				isNewMesssage = true
			}
			//fmt.Println(original_response)
			if original_response.Error != nil ||
				((original_response.Message.Author.Name == "assistant" ||
					original_response.Message.Author.Role == "assistant") &&
					original_response.Message.Status == "finished_successfully" &&
					original_response.Message.Metadata.MessageType == "continue" &&
					original_response.Message.EndTurn == false) {
				// Try to continue
				translated_request.Action = "continue"
				if translated_request.Model == "gpt-4-browsing" {
					translated_request.Model = "gpt-4"
				}
				translated_request.ParentMessageID = lastMessageId
				translated_request.ConversationID = original_response.ConversationID
				translated_request.Messages = nil
				response, err = chatgpt.SendRequest(translated_request, token)
				if err != nil {
					fmt.Println(err.Error())
					c.JSON(500, gin.H{
						"error": "error sending request" + err.Error(),
					})
					return
				}
				defer response.Body.Close()
				if response.StatusCode != 200 {
					// Try read response body as JSON
					var error_response map[string]interface{}
					err = json.NewDecoder(response.Body).Decode(&error_response)
					if err != nil {
						c.JSON(500, gin.H{"error": gin.H{
							"message": "Unknown error: " + err.Error(),
							"type":    "internal_server_error",
							"param":   nil,
							"code":    "500",
						}})
						return
					}
					c.JSON(response.StatusCode, gin.H{"error": gin.H{
						"message": error_response["detail"],
						"type":    response.Status,
						"param":   nil,
						"code":    "error",
					}})
					return
				}
				//bytes, err := ioutil.ReadAll(response.Body)
				//fmt.Println(string(bytes))

				// Create a bufio.Reader from the response body
				reader = bufio.NewReader(response.Body)
				continue
			}
			if original_response.Message.Author.Role == "tool" &&
				original_response.Message.Author.Name == "browser" &&
				original_response.Message.Status == "finished_successfully" &&
				len(original_response.Message.Metadata.CiteMetadata.MetaDataList) > 0 {
				last_browser_metadata.CiteMetadata.MetaDataList = []responses.MetaDataListItem{}
				for _, item := range original_response.Message.Metadata.CiteMetadata.MetaDataList {
					last_browser_metadata.CiteMetadata.MetaDataList = append(last_browser_metadata.CiteMetadata.MetaDataList, item)
				}
			}
			if original_response.Message.Author.Role != "assistant" {
				continue
			}
			if original_response.Message.Metadata.Timestamp == "absolute" {
				continue
			}
			translated_response := responses.ChatCompletionChunk{}

			tmp_fulltext := ""

			switch original_response.Message.Content.ContentType {
			case "text":
				if len(original_response.Message.Content.Parts) == 0 || original_response.Message.Content.Parts[0] == "" {
					continue
				}
				tmp_fulltext = original_response.Message.Content.Parts[0]
				original_response.Message.Content.Parts[0] = strings.ReplaceAll(original_response.Message.Content.Parts[0], fulltext, "")
				text := original_response.Message.Content.Parts[0]
				if original_request.Model == "gpt-4-browsing" {
					if isNewMesssage {
						text = "///text_message\n" + text
					}
					if original_response.Message.Status == "finished_successfully" {
						if len(original_response.Message.Metadata.Citations) > 0 {
							last_browser_metadata.Citations = []responses.Citation{}
							for _, item := range original_response.Message.Metadata.Citations {
								last_browser_metadata.Citations = append(last_browser_metadata.Citations, item)
							}
						}

						metadataString, _ := json.Marshal(last_browser_metadata)
						text += "\n" + "%%%TEXT_METADATA:" + string(metadataString) + "%%%\n" + `\\\` + "\n"
					}
				}
				translated_response = responses.NewChatCompletionChunk(original_request.Model, text)
			case "code":
				tmp_fulltext = original_response.Message.Content.Text

				original_response.Message.Content.Text = strings.ReplaceAll(original_response.Message.Content.Text, fulltext, "")

				text := original_response.Message.Content.Text
				if original_request.Model == "gpt-4-browsing" {
					if isNewMesssage {
						metadataString, _ := json.Marshal(last_browser_metadata)
						text = "///code_message\n%%%METADATA:" + string(metadataString) + "%%%" + text + "\n"
					}
					if original_response.Message.Status == "finished_successfully" {
						text += "\n" + `\\\` + "\n"
					}
				}
				translated_response = responses.NewChatCompletionChunk(original_request.Model, text)
			}

			// Stream the response to the client
			response_string := translated_response.String()
			if original_request.Stream {
				//fmt.Println(string(response_string) + "\n")
				_, err = c.Writer.WriteString("data: " + string(response_string) + "\n\n")
				if err != nil {
					return
				}
			}

			// Flush the response writer buffer to ensure that the client receives each line as it's written
			c.Writer.Flush()
			fulltext = tmp_fulltext
			lastMessageId = original_response.Message.ID
		}
	}

}
