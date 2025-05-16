package handlers

import (
	"database/sql"

	"example.com/weather-report/config"
	"example.com/weather-report/restapi/operations/subscription"
	"example.com/weather-report/restapi/operations/weather"
	// "example.com/weather-report/storage"
	"example.com/weather-report/weather-api"
	"github.com/go-openapi/runtime/middleware"
)

func subscribeHandler(db *sql.DB, params subscription.SubscribeParams) middleware.Responder {
	return middleware.NotImplemented("")
}

func unsubscribeHandler(db *sql.DB, params subscription.UnsubscribeParams) middleware.Responder {
	return middleware.NotImplemented("")
}

func confirmSubscriptionHandler(db *sql.DB, params subscription.ConfirmSubscriptionParams) middleware.Responder {
	return middleware.NotImplemented("")
}

func getWeatherHandler(conf *config.WeatherApiConfig,params weather.GetWeatherParams) middleware.Responder {
	response, err := weatherApi.GetCurrentWeather(params.City, conf)
	if err != nil {
		return weather.NewGetWeatherNotFound()
	}
	handlerRespone := weather.NewGetWeatherOK()
	handlerRespone.SetPayload(&weather.GetWeatherOKBody{
		Description: response.Description,
		Humidity: float64(response.Humidity),
		Temperature: response.TempC,
	})
	return handlerRespone
}