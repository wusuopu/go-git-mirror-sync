package config

import (
	"os"
)

type IConfig map[string]interface{}

var Config IConfig

func Load() IConfig {
	if Config == nil {
		PORT := os.Getenv("PORT")
		GO_ENV := os.Getenv("GO_ENV")
		if PORT == "" {
			PORT = "80"
		}
		if GO_ENV == "" {
			GO_ENV = "development"
		}

		Config = IConfig{
			"PORT": PORT,
			"GO_ENV": GO_ENV,
		}
	}
	return Config
}
