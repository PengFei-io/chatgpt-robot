package intention

import (
	"chatgpt-robot/ai"
	"chatgpt-robot/utils/joke"
	"chatgpt-robot/utils/weather"
	"context"
	"fmt"
	"strings"
	"time"
)

// IntentHandle 意图识别,判断天气.
func IntentHandle(fromUserName string, content string) string {
	t0 := time.Now()
	defer func() {
		t1 := time.Now()
		fmt.Printf("消耗时间:%dms\n", t1.Sub(t0).Milliseconds())
	}()
	if strings.Contains(content, "@天气") {
		fmt.Println("hit_intent_天气")
		return weather.GetWeather(content)
	} else if strings.Contains(content, "@笑话") {
		fmt.Println("hit_intent_笑话")
		return joke.GetJokeList(5)
	}
	fmt.Println("hit_intent_chat")
	ctx := context.Background()
	return ai.GetSessionOpenAITextReply(ctx, content, fromUserName)
}
