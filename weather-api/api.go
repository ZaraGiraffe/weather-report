package weatherApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"example.com/weather-report/config"
)

func CreateUrl(conf *config.WeatherApiConfig, city string) string {
	return fmt.Sprintf("%s/current.json?key=%s&q=%s", conf.Url, conf.ApiKey, city) 
}

type weatherApiJsonResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		Humidity int `json:"humidity"`
		Condition struct {
			Description string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

type WeatherResponse struct {
	TempC float64
	Humidity int
	Description string
}

type errorJsonResponse struct {
    Error struct {
        Code int `json:"code"`
        Message string `json:"message"`
    } `json:"error"`
}

func (w *weatherApiJsonResponse) ToWeatherResponse() WeatherResponse {
	return WeatherResponse{
		TempC: w.Current.TempC,
		Humidity: w.Current.Humidity,
		Description: w.Current.Condition.Description,
	}
}

var ErrBadInput error = errors.New("bad input to the request")

func GetCurrentWeather(city string, conf *config.WeatherApiConfig) (*WeatherResponse, error) {
    requestUrl := CreateUrl(conf, city)

    req, err := http.Get(requestUrl)
    if err != nil {
        log.Println("Error creating request:", err)
        return nil, err
    }
    defer req.Body.Close()

    bodyBytes, err := io.ReadAll(req.Body)
    if err != nil {
        log.Println("Error reading response body:", err)
        return nil, err
    }

    var weatherJsonResponse weatherApiJsonResponse
    var errorJsonResponse errorJsonResponse
    err = json.Unmarshal(bodyBytes, &errorJsonResponse)

    if (err == nil && errorJsonResponse.Error.Code != 0) || req.StatusCode != http.StatusOK {
        log.Println("Error input in request")
        return nil, ErrBadInput
    }

    err = json.Unmarshal(bodyBytes, &weatherJsonResponse)
    if err != nil {
        log.Fatalf("ERROR: in get current weather: %v", err)
    }

    weatherResponse := weatherJsonResponse.ToWeatherResponse()
    return &weatherResponse, nil
}