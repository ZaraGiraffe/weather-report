package processor

import (
	"database/sql"
	"log"
	"time"

	"example.com/weather-report/config"
	"example.com/weather-report/emails"
	"example.com/weather-report/storage"
	"example.com/weather-report/weather-api"
)

var SLEEP_TIME_MINUTES int = 1

type Processor struct {
	conf *config.Config
	db   *sql.DB
}

func NewProcessor(conf *config.Config, db *sql.DB) *Processor {
	return &Processor{
		conf: conf,
		db:   db,
	}
}

func (p *Processor) Run() {
	for {
		frequencies := []int{config.HOURLY_FREQUENCY, config.DAILY_FREQUENCY}
		for _, frequency := range frequencies {
			timeConstraint := time.Now().Add(-time.Duration(frequency) * time.Minute).Unix()
			rows := storage.GetAllSubscriptionsWithTimeConstraint(p.db, timeConstraint, frequency)
			for _, row := range rows {
				if row.Status == config.CONFIRMED_STATUS {
					weatherReport, err := weatherApi.GetCurrentWeather(row.City, &p.conf.WeatherApiConfig)
					if err != nil {
						log.Fatalf("ERROR: get current weather query failed: %v", err)
					}
					emails.SendWeatherReportEmail(row.Email, weatherReport, p.conf)
					storage.UpdateSubscriptionLastSent(p.db, row.Token, time.Now().Unix())
				}
			}
		}
		time.Sleep(time.Duration(SLEEP_TIME_MINUTES) * time.Minute)
	}
}
