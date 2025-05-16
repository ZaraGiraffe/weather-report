package handlers

import (
	"database/sql"
	"github.com/go-openapi/runtime/middleware"
	"example.com/weather-report/restapi/operations/subscription"
	"example.com/weather-report/restapi/operations/weather"
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

func getWeatherHandler(db *sql.DB, params weather.GetWeatherParams) middleware.Responder {
	return middleware.NotImplemented("")
}