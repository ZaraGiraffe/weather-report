// This package is responsible for handling the requests from the client
// It contains the handlers for the API endpoints
package handlers

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"example.com/weather-report/config"
	"example.com/weather-report/emails"
	"example.com/weather-report/restapi/operations/subscription"
	"example.com/weather-report/restapi/operations/weather"
	"example.com/weather-report/storage"
	weatherApi "example.com/weather-report/weather-api"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
)

func subscribeHandler(conf *config.Config, db *sql.DB, params subscription.SubscribeParams) middleware.Responder {
	_, err := weatherApi.GetCurrentWeather(params.City, &conf.WeatherApiConfig)
	if err != nil {
		log.Println("INFO: Error bad city name")
		return subscription.NewSubscribeBadRequest()
	}

	var frequencyType int
	switch params.Frequency {
	case "hourly":
		frequencyType = config.HOURLY_FREQUENCY
	case "daily":
		frequencyType = config.DAILY_FREQUENCY
	default:
		log.Println("INFO: Error bad frequency param")
		return subscription.NewSubscribeBadRequest()
	}

	_, err = storage.GetSubscriptionByEmail(db, params.Email)
	if !errors.Is(err, sql.ErrNoRows) {
		log.Println("INFO: Email already exists")
		return subscription.NewSubscribeConflict()
	}

	token := uuid.New().String()
	new_subscription := &storage.Subscription{
		Email:         params.Email,
		City:          params.City,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
		FrequencyType: frequencyType,
		Token:         token,
		Status:        config.PENDING_STATUS,
	}
	err = storage.InsertSubscriptionQuery(db, new_subscription)
	if err != nil {
		log.Printf("ERROR: problem with inserting new subscription: %v", err)
		return subscription.NewSubscribeBadRequest()
	}

	err = emails.SendConfirmationEmail(params.Email, token, conf)
	if err != nil {
		log.Printf("INFO: Email has failed to be sent, seems like email is invalid: %v", params.Email)
		err := storage.DeleteSubscriptionByToken(db, token)
		if err != nil {
			log.Printf("ERROR: problem with deleting subscription with bad email: %v", err)
		}
		return subscription.NewSubscribeBadRequest()
	}

	log.Printf("INFO: Subscribed to token %v", token)
	return subscription.NewSubscribeOK()
}

func unsubscribeHandler(db *sql.DB, params subscription.UnsubscribeParams) middleware.Responder {
	_, err := storage.GetSubscriptionByToken(db, params.Token)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("INFO: Subscription with token %v not found", params.Token)
		return subscription.NewUnsubscribeNotFound()
	}

	err = storage.DeleteSubscriptionByToken(db, params.Token)
	if err != nil {
		log.Printf("ERROR: problem with deleting subscription with bad token: %v", err)
		return subscription.NewUnsubscribeBadRequest()
	}

	log.Printf("INFO: Unsubscribed from token %v", params.Token)
	return subscription.NewUnsubscribeOK()
}

func confirmSubscriptionHandler(db *sql.DB, params subscription.ConfirmSubscriptionParams) middleware.Responder {
	_, err := storage.GetSubscriptionByToken(db, params.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("INFO: Subscription with token %v not found", params.Token)
			return subscription.NewConfirmSubscriptionNotFound()
		} else {
			log.Printf("ERROR: problem with getting subscription by token: %v", err)
			return subscription.NewConfirmSubscriptionBadRequest()
		}
	}

	err = storage.UpdateSubscriptionStatus(db, params.Token, config.CONFIRMED_STATUS)
	if err != nil {
		log.Printf("ERROR: problem with updating subscription status: %v", err)
		return subscription.NewConfirmSubscriptionBadRequest()
	}

	log.Printf("INFO: Confirmed subscription with token %v", params.Token)
	return subscription.NewConfirmSubscriptionOK()
}

func getWeatherHandler(conf *config.WeatherApiConfig, params weather.GetWeatherParams) middleware.Responder {
	response, err := weatherApi.GetCurrentWeather(params.City, conf)
	if err != nil {
		if errors.Is(err, weatherApi.ErrBadInput) {
			log.Printf("INFO: City %v not found", params.City)
			return weather.NewGetWeatherNotFound()
		} else {
			log.Printf("ERROR: problem with getting weather: %v", err)
			return weather.NewGetWeatherBadRequest()
		}
	}

	handlerRespone := weather.NewGetWeatherOK()
	handlerRespone.SetPayload(&weather.GetWeatherOKBody{
		Description: response.Description,
		Humidity:    float64(response.Humidity),
		Temperature: response.TempC,
	})

	log.Printf("INFO: Weather for city %v fetched successfully", params.City)
	return handlerRespone
}
