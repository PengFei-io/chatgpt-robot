package main

import (
	config2 "chatgpt-robot/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config2.LoadConfig()
	config2.ShowConfig()

	r := gin.New()
	// 微信公众号接口
	r.GET("/wechat", wxCheckSign)
	r.POST("/wechat", wxChatMessage)

	// 捷径接口
	//r.POST("/chatgpt/api/completions", completions)

	config := config2.GetConfig()
	if err := r.Run(fmt.Sprintf(":%v", config.Port)); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
