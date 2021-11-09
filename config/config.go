package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	BOT_TOKEN string
	BOT_API   string
	BOT_URL   string
	PORT      string
	IS_LOCAL  bool
}

func New() *Config {
	conf := Config{
		BOT_TOKEN: getRequiredEnv("BOT_TOKEN"),
		BOT_API:   getRequiredEnv("BOT_API"),
		PORT:      getEnv("PORT", "8080"),
	}
	conf.BOT_URL = conf.BOT_API + conf.BOT_TOKEN
	isLocal, err := strconv.ParseBool(getEnv("IS_LOCAL", "true"))
	if err != nil {
		isLocal = true
	}
	conf.IS_LOCAL = isLocal

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
