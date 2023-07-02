package weather

import (
	"testing"
)

func TestGetCityWeather(t *testing.T) {
	content := "上海天气如何"
	weatherContent := GetWeather(content)
	println(weatherContent)
}
