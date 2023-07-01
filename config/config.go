package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync/atomic"

	"github.com/spf13/viper"
)

// Config 配置内容
type Config struct {
	Key        string
	Port       int32
	Token      string
	Timeout    int32
	SessionTTL int32 `yaml:"session_ttl"`
}

var Prompt string

var gConfig atomic.Value

// LoadConfig 加载配置
func LoadConfig() {
	// 读取配置文件
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s", err))
	}

	// 解析配置文件
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %s", err))
	}
	gConfig.Store(config)
	prompt, err := ioutil.ReadFile("prompt.txt")
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `prompt.txt` 的配置, 错误信息: %+v\n", err)
	}
	Prompt = string(prompt)
}

// GetConfig 获取配置
func GetConfig() Config {
	return gConfig.Load().(Config)
}

// ShowConfig 展示配置信息
func ShowConfig() {
	config := GetConfig()
	log.Println("====================================================")
	log.Printf("OpenAIKey=%v\n", config.Key)
	log.Printf("Token =%v\n", config.Token)
	log.Printf("port=%v\n", config.Port)
	log.Printf("timeout=%v\n", config.Timeout)
	log.Println("====================================================")
}