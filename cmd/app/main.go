package main

import (
	"example.com/weather-report/emails/processor"
	"example.com/weather-report/config"
	"example.com/weather-report/storage"
)
	

func main() {
	config.LoadEnv()

	processor := processor.NewProcessor(config.GetConfig(), storage.NewStorageConnection())
	go processor.Run()

	runServer()
}