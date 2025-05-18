package weatherApi

import (
	"testing"

	"example.com/weather-report/config"
	"github.com/stretchr/testify/assert"
)

func Test_MakeWeatherApiRequest(t *testing.T) {
	conf := config.GetConfig("../test.config.json")
	report, err := GetCurrentWeather("London", &conf.WeatherApiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, report)
}
