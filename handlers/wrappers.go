package handlers

import (
	"database/sql"
	
	"github.com/go-openapi/runtime/middleware"
	"example.com/weather-report/restapi/operations/subscription"
	"example.com/weather-report/restapi/operations/weather"
	"example.com/weather-report/config"
)

func SubscribeHandlerWrapper(storage *sql.DB) func(subscription.SubscribeParams) middleware.Responder {
	return func(params subscription.SubscribeParams) middleware.Responder {
		return subscribeHandler(storage, params)
	}
}

func UnsubscribeHandlerWrapper(storage *sql.DB) func(subscription.UnsubscribeParams) middleware.Responder {
	return func(params subscription.UnsubscribeParams) middleware.Responder {
		return unsubscribeHandler(storage, params)
	}
}

func ConfirmSubscriprionHandlerWrapper(storage *sql.DB) func(subscription.ConfirmSubscriptionParams) middleware.Responder {
	return func(params subscription.ConfirmSubscriptionParams) middleware.Responder {
		return confirmSubscriptionHandler(storage, params)
	}
}

func GetWeatherHandlerWrapper(conf *config.WeatherApiConfig) func(weather.GetWeatherParams) middleware.Responder {
	return func(params weather.GetWeatherParams) middleware.Responder {
		return getWeatherHandler(conf, params)
	}
}