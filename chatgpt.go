package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang/groupcache"
	"github.com/sashabaranov/go-openai"
)

var cache *groupcache.Group

func init() {
	// 初始化缓存，最多缓存128M数据
	cache = groupcache.NewGroup("chatgpt", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			// 从后端数据源获取数据
			elems := strings.Split(key, "______")
			value, err := chatCompletion(elems[0], elems[1])
			if err != nil {
				log.Printf("chatCompletion Err=%v\n", err)
				return err
			}
			dest.SetString(value)
			return nil
		}))
}

// 构建缓存key
func getKey(key, content string) string {
	return fmt.Sprintf("%v______%v", key, content)
}

// GetChatData 获取OPenAI聊天数据
func GetChatData(appKey, content string) string {
	key := getKey(appKey, content)
	var value string
	err := cache.Get(nil, key, groupcache.StringSink(&value))
	if err != nil {
		log.Fatalf("cache.Get%v\n", err.Error())
		return "你的问题太复杂了,你可以再问我一遍吗."
	}
	return value
}

// 从openAI接口取数据
func chatCompletion(key, content string) (string, error) {

	cfg := openai.DefaultConfig(key)
	cfg.HTTPClient = &http.Client{Timeout: 4500 * time.Millisecond} // <= add a custom http client
	client := openai.NewClientWithConfig(cfg)
	rsp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo0301,
			Stream:    false,
			MaxTokens: 2048, // 限制最大返回token，提速
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "你的问题太复杂了,我可能要再想一会.", err
	}
	if len(rsp.Choices) == 0 {
		return "你的问题太难了,我可能不会了,太难了.", fmt.Errorf("Get empty response")
	}
	return rsp.Choices[0].Message.Content, nil
}
