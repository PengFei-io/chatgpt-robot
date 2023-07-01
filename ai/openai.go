package ai

import (
	"chatgpt-robot/config"
	"chatgpt-robot/utils"
	"context"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	ErrTips = "抱歉，出错了，请稍后重试~"
	GptApi  = "https://api.openai.com/v1"
)

func getOpenAIClient(apiKey string) *openai.Client {
	var c openai.ClientConfig
	c = openai.DefaultConfig(apiKey)
	c.BaseURL = GptApi
	return openai.NewClientWithConfig(c)
}

func CreateChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) string {
	getConfig := config.GetConfig()
	apiKey := getConfig.Key
	client := getOpenAIClient(apiKey)
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    getConfig.Model,
		Messages: messages,
	})
	if err != nil {
		log.Printf("openAIClient.CreateChatCompletion err=%+v\n", err)
		return ErrTips
	}
	if len(resp.Choices) == 0 {
		log.Printf("resp is err, resp=%s", utils.MarshalAnyToString(resp))
		return ErrTips
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content)
}
