package intention

import (
	"chatgpt-robot/ai"
	"chatgpt-robot/utils/weather"
	"context"
	"strings"
)

// IntentHandle 意图识别,判断天气.
func IntentHandle(fromUserName string, content string) string {
	if strings.Contains(content, "@天气") {
		return weather.GetWeather(content)
	}
	ctx := context.Background()
	return ai.GetSessionOpenAITextReply(ctx, content, fromUserName)
}
