package weather

import (
	"strings"
	"testing"
)

func TestGetCityWeather(t *testing.T) {
	content := "#天气沭阳"
	if strings.Contains(content, "#天气") {
		weatherContent := GetWeather(content)
		println(weatherContent)
	}
}
