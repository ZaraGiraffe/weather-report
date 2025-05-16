package config


import (
    "log"
    "os"
	"io"
	"encoding/json"
	"github.com/joho/godotenv"
)

type StorageConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type WeatherApiConfig struct {
	ApiKey string `json:"api-key"`
	Url string `json:"url"`
}

type Config struct {
	AppPassword string `json:"app-password"`
	WeatherApiConfig WeatherApiConfig `json:"weather-api"`
	StorageConfig StorageConfig `json:"storage-config"`
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getConfig(configPath string) *Config {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal("Error opening config file: ", err)
	}
	defer file.Close()

	fileData, _ := io.ReadAll(file)

	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		log.Fatal("Error unmarshalling config file: ", err)
	}
	
	return &config
}

func GetConfig() *Config {
	loadEnv()
	configPath, present := os.LookupEnv("CONFIG")
	if !present {
		log.Fatal("CONFIG not found in .env file")
	}
	return getConfig(configPath)
}