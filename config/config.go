package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	TWOPOI_BOT_TOKEN string
	TWOPOI_BOT_API   string
	TWOPOI_BOT_URL   string
	TWOPOI_HOST      string
	TWOPOI_PORT      string
	TWOPOI_IS_LOCAL  bool
}

func New() *Config {
	conf := Config{
		TWOPOI_BOT_TOKEN: getRequiredEnv("TWOPOI_BOT_TOKEN"),
		TWOPOI_BOT_API:   getRequiredEnv("TWOPOI_BOT_API"),
		TWOPOI_HOST:      getRequiredEnv("TWOPOI_HOST"),
		TWOPOI_PORT:      getEnv("TWOPOI_PORT", "8080"),
	}
	conf.TWOPOI_BOT_URL = conf.TWOPOI_BOT_API + conf.TWOPOI_BOT_TOKEN
	isLocal, err := strconv.ParseBool(getEnv("TWOPOI_IS_LOCAL", "true"))
	if err != nil {
		isLocal = true
	}
	conf.TWOPOI_IS_LOCAL = isLocal

	return &conf
}

func getRequiredEnv(key string) string {
	value, isFound := os.LookupEnv(key)
	if !isFound {
		log.Fatalf("Environment variable \"%s\" not found.\n", key)
	}
	return value
}

func getEnv(key string, defaultValue string) string {
	value, isFound := os.LookupEnv(key)
	if isFound {
		return value
	}
	return defaultValue
}
