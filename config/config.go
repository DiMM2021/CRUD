package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

var log = logrus.New()

func LoadConfig() *Config {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel)

	config := &Config{
		DBHost:     getEnv("DB_HOST", "postgres_container"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "qwerty"),
		DBName:     getEnv("DB_NAME", "cruddb"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.WithFields(logrus.Fields{
			"key":           key,
			"default_value": defaultValue,
		}).Warn("Environment variable not set, using default value")
		return defaultValue
	}
	return value
}
