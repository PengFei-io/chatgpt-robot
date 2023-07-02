package weather

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Weather struct {
	City     string
	CityCode string
	Country  string
	Province string
}

var citySlice []Weather

const weatherApi = "https://devapi.qweather.com/v7/weather/now?key=9b83895fef154e9bb02cb567c0817b4e&location=%s&lang=zh"

func init() {
	// 打开CSV文件
	file, err := os.Open("./static/China-City-List-latest.csv")
	//file, err := os.Open("/Users/pf/dev/IdeaProjects/go/chatgpt-robot/static/China-City-List-latest.csv")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 创建CSV Reader
	reader := csv.NewReader(file)

	// 读取CSV文件中的数据
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	// 打印CSV文件中的数据
	citySlice = make([]Weather, 0)
	for _, record := range records {
		w := Weather{record[2], record[0], record[7], record[9]}
		citySlice = append(citySlice, w)
	}
}

// GetCityCode 获取城市的ID
func GetCityCode(city string) []Weather {
	weatherSlice := make([]Weather, 0)
	for _, v := range citySlice {
		if strings.Contains(city, v.City) {
			weatherSlice = append(weatherSlice, v)
		}
	}
	return weatherSlice
}

// GetCityWeather 获取城市的天气
func GetCityWeather(city string) []map[string]map[string]any {
	// 创建一个http.Request对象
	weatherSlice := GetCityCode(city)
	weatherRetSlice := make([]map[string]map[string]any, 0)
	for _, w := range weatherSlice {
		m := make(map[string]map[string]any)
		req, err := http.NewRequest("GET", fmt.Sprintf(weatherApi, w.CityCode), nil)
		if err != nil {
			log.Println(err)
		}
		// 添加请求头
		req.Header.Set("User-Agent", "Mozilla/5.0")

		// 创建一个http.Client对象
		client := &http.Client{}

		// 发送HTTP请求
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		// 读取响应内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		var weatherMap map[string]any
		err = json.Unmarshal(body, &weatherMap)
		if err != nil {
			log.Println(err)
		}
		m[w.Country+" "+w.Province+" "+w.City] = weatherMap
		weatherRetSlice = append(weatherRetSlice, m)
	}
	return weatherRetSlice
}

// formatWeatherPrompts 格式化天气的提示语
func formatWeatherPrompts(hitCity string, weatherMap map[string]any) string {
	updateTime := weatherMap["updateTime"]
	fxLink := weatherMap["fxLink"]
	now := weatherMap["now"]
	nowMap := now.(map[string]interface{})

	nowStr := fmt.Sprintf(
		"%s,天气状况:%s,温度:%s°C.体感温度:%s°C.风向:%s.风力等级:%s.能见度:%s.云量:%s",
		hitCity,
		nowMap["text"],
		nowMap["temp"],
		nowMap["feelsLike"],
		nowMap["windDir"],
		nowMap["windScale"],
		nowMap["vis"],
		nowMap["cloud"],
	)
	return fmt.Sprintf("%s,更新时间:%s,详细信息可以访问:%s", nowStr, FormatDateStr(fmt.Sprintf("%v", updateTime)), fxLink)
}

// GetWeather 获取城市天气
func GetWeather(city string) string {
	weatherSlice := GetCityWeather(city)
	var weatherPrompts []string
	for _, s := range weatherSlice {
		for hitCity, weatherMap := range s {
			weatherPrompts = append(weatherPrompts, formatWeatherPrompts(hitCity, weatherMap))
		}
	}
	return strings.Join(weatherPrompts, "\n")
	//return formatWeatherPrompts(hitCity, weatherMap)
}

func FormatDateStr(input string) string {
	// 将时间字符串解析为time.Time类型
	t, err := time.Parse("2006-01-02T15", input[:13])
	if err != nil {
		log.Println(err)
	}
	// 将time.Time类型的时间格式化为指定的日期和时间格式
	formatted := t.Format("2006-01-02 15:04:05")
	// 打印格式化后的时间
	return formatted
}
