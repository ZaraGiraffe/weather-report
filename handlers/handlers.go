package handlers

import (
	"database/sql"
	"example.com/weather-report/storage"
	"example.com/weather-report/config"
	"example.com/weather-report/restapi/operations/subscription"
	"example.com/weather-report/restapi/operations/weather"
	"example.com/weather-report/weather-api"
	"example.com/weather-report/emails"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"time"
	"errors"
	"log"
)

func subscribeHandler(conf *config.Config, db *sql.DB, params subscription.SubscribeParams) middleware.Responder {
	_, err := weatherApi.GetCurrentWeather(params.City, &conf.WeatherApiConfig)
	if err != nil {
		log.Println("Error bad city name")
		return subscription.NewSubscribeBadRequest()
	}	
	var frequencyType int
	switch params.Frequency {
	case "hourly":
		frequencyType = config.HOURLY_FREQUENCY
	case "daily":
		frequencyType = config.DAILY_FREQUENCY
	default:
		log.Println("Error bad frequency param")
		return subscription.NewSubscribeBadRequest()
	}

	_, err = storage.GetSubscriptionByEmail(db, params.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		log.Println("Email already exists")
		return subscription.NewSubscribeConflict()
	}

	token := uuid.New().String()
	new_subscription := &storage.Subscription{
		Email: params.Email,
		City: params.City,
		Created_at: time.Now().Unix(),
		Updated_at: time.Now().Unix(),
		Frequency_type: frequencyType,
		Token: token,
		Status: config.PENDING_STATUS,
	}
	err = storage.InsertDubscriptionQuery(db, new_subscription)
	if err != nil {
		log.Fatal(err)
	}

	err = emails.SendConfirmationEmail(params.Email, token, conf)
	if err != nil {
		log.Println("Email has failed to be sent")
		err := storage.DeleteSubscription(db, token)
		if err != nil {
			log.Fatalf("ERROR: delete subscription query failed: %v", err)
		}
		return subscription.NewSubscribeBadRequest()
	}

	return subscription.NewSubscribeOK()
}

func unsubscribeHandler(db *sql.DB, params subscription.UnsubscribeParams) middleware.Responder {
	_, err := storage.GetSubscriptionByToken(db, params.Token)
	if errors.Is(err, sql.ErrNoRows) {
		return subscription.NewUnsubscribeNotFound()
	}
	storage.DeleteSubscription(db, params.Token)
	return subscription.NewUnsubscribeOK()
}

func confirmSubscriptionHandler(db *sql.DB, params subscription.ConfirmSubscriptionParams) middleware.Responder {
	_, err := storage.GetSubscriptionByToken(db, params.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return subscription.NewConfirmSubscriptionNotFound()
		}
	}
	storage.UpdateSubscriptionStatus(db, params.Token, config.CONFIRMED_STATUS)
	return subscription.NewConfirmSubscriptionOK()
}

func getWeatherHandler(conf *config.WeatherApiConfig,params weather.GetWeatherParams) middleware.Responder {
	response, err := weatherApi.GetCurrentWeather(params.City, conf)
	if errors.Is(err, weatherApi.ErrBadInput) {
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