package weatherApi

import (
	"example.com/weather-report/config"
	"net/http"
	"fmt"
	"encoding/json"
	"io"
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

func (w *weatherApiJsonResponse) ToWeatherResponse() WeatherResponse {
	return WeatherResponse{
		TempC: w.Current.TempC,
		Humidity: w.Current.Humidity,
		Description: w.Current.Condition.Description,
	}
}

func GetCurrentWeather(city string, conf *config.WeatherApiConfig) (*WeatherResponse, error) {
    requestUrl := CreateUrl(conf, city)

    req, err := http.Get(requestUrl)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return nil, err
    }
    defer req.Body.Close()

    if req.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %d, status text: %s", req.StatusCode, req.Status)
    }

    bodyBytes, err := io.ReadAll(req.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        return nil, err
    }

    var weatherJsonResponse weatherApiJsonResponse
    err = json.Unmarshal(bodyBytes, &weatherJsonResponse)
    if err != nil {
        fmt.Println("Error unmarshalling response:", err)
        return nil, err
    }

    weatherResponse := weatherJsonResponse.ToWeatherResponse()
    return &weatherResponse, nil
}