// This package is responsible for sending regular weather report emails to the users
// It always runs in the background and sends emails to the users based on the frequency of the subscriptions
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

// The time to sleep between each iteration of the loop
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
		log.Println("INFO: Starting processor cycle")

		frequencies := []int{config.HOURLY_FREQUENCY, config.DAILY_FREQUENCY}
		for _, frequency := range frequencies {
			timeConstraint := time.Now().Add(-time.Duration(frequency) * time.Minute).Unix()
			rows, err := storage.GetAllSubscriptionsWithTimeConstraint(p.db, timeConstraint, frequency)
			if err != nil {
				log.Printf("ERROR: problem with getting all subscriptions with time constraint: %v", err)
				continue
			}

			for _, row := range rows {
				if row.Status == config.CONFIRMED_STATUS {
					weatherReport, err := weatherApi.GetCurrentWeather(row.City, &p.conf.WeatherApiConfig)
					if err != nil {
						log.Printf("ERROR: get current weather query failed, maybe it was deleted: %v", err)
						continue
					}
					emails.SendWeatherReportEmail(row.Email, weatherReport, p.conf)
					err = storage.UpdateSubscriptionLastSent(p.db, row.Token, time.Now().Unix())
					if err != nil {
						log.Printf("ERROR: problem with updating subscription last sent: %v", err)
					}
				}
			}
		}

		time.Sleep(time.Duration(SLEEP_TIME_MINUTES) * time.Minute)
	}
}
