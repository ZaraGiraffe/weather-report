// This package is responsible for loading the configuration from the config file
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type WeatherApiConfig struct {
	ApiKey      string `json:"api-key"`
	Url         string `json:"url"`
	AppPassword string `json:"app-password"`
}

type AdminEmailConfig struct {
	Email       string `json:"email"`
	AppPassword string `json:"app-password"`
}

type SmtpServerConfig struct {
	Host    string `json:"host"`
	TlsPort string `json:"tls-port"`
}

type ServerConfig struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	TlsEnabled bool   `json:"tls-enabled"`
}

func (s *ServerConfig) GetAddress() string {
	protocol := "http"
	if s.TlsEnabled {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, s.Host, s.Port)
}

type Config struct {
	SmtpServerConfig SmtpServerConfig `json:"smtp-server-config"`
	AdminEmailConfig AdminEmailConfig `json:"admin-email-config"`
	WeatherApiConfig WeatherApiConfig `json:"weather-api-config"`
	StorageConfig    StorageConfig    `json:"storage-config"`
	ServerConfig     ServerConfig     `json:"server-config"`
}

func LoadEnv(envPath ...string) {
	var err error
	if len(envPath) > 0 {
		err = godotenv.Load(envPath[0])
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

func GetConfig(optConfigPath ...string) *Config {
	configPath := ""
	if len(optConfigPath) > 0 {
		configPath = optConfigPath[0]
	} else {
		envConfigPath, present := os.LookupEnv("CONFIG")
		if !present {
			log.Fatal("CONFIG not found in .env file")
		}
		configPath = envConfigPath
	}
	return getConfig(configPath)
}
