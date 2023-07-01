package main

import (
	config2 "chatgpt-robot/config"
	"chatgpt-robot/utils/intention"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// iPhone捷径请求消息
type Request struct {
	Prompt  string `json:"prompt"`
	Content string `json:"content"`
}

// 定义微信消息结构体
type Message struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        int64    `xml:"MsgId"`
}

// 处理iPhone捷径请求
func completions(c *gin.Context) {
	//key := c.GetHeader("Authorization")

	request := &Request{}
	c.BindJSON(request)

	log.Println("request=", request.Content)
	content := intention.IntentHandle(request.Content, "")

	c.JSON(http.StatusOK, gin.H{
		"content": content,
	})
	log.Println("response=", content)
	return
}

func wxChatMessage(c *gin.Context) {
	// 读取POST请求中的消息内容
	msg := &Message{}
	err := xml.NewDecoder(c.Request.Body).Decode(msg)
	if err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}
	marshalMsg, _ := json.Marshal(msg)
	fmt.Println(string(marshalMsg))
	// 构造回复消息
	resp := &Message{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   msg.CreateTime,
		MsgType:      "text",
		Content:      "",
	}
	//只处理text类型的消息
	if "text" != msg.MsgType {
		//resp.Content = fmt.Sprintf("%v类型目前不支持,抱歉.你可以发文本给我吗?\n", msg.Content)
		//// 返回响应
		//c.Writer.Header().Set("Content-Type", "application/xml")
		//respXML, _ := xml.Marshal(resp)
		//c.Writer.Write(respXML)
		return
	}

	// Get response from OpenAI
	log.Printf("request=%v\n", msg.Content)
	content := intention.IntentHandle(msg.FromUserName, msg.Content)

	log.Println("response=", content)
	resp.Content = content
	respXML, err := xml.Marshal(resp)
	if err != nil {
		fmt.Println("Error encoding message:", err)
		return
	}

	// 返回响应
	c.Writer.Header().Set("Content-Type", "application/xml")
	c.Writer.Write(respXML)
}

func wxCheckSign(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")

	// 将token、timestamp、nonce三个参数进行字典序排序
	config := config2.GetConfig()
	strs := sort.StringSlice{config.Token, timestamp, nonce}
	strs.Sort()

	// 将三个参数字符串拼接成一个字符串进行sha1加密
	str := strings.Join(strs, "")
	h := sha1.New()
	io.WriteString(h, str)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	// 将加密后的字符串与signature进行对比，判断该请求是否来自微信服务器
	if hash == signature {
		c.String(http.StatusOK, echostr)
	} else {
		c.String(http.StatusBadRequest, "Invalid signature")
	}
}
