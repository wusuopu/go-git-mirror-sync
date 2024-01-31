package config

import (
	"os"
)

type IConfig map[string]interface{}

var Config IConfig

func Load() IConfig {
	if Config == nil {
		PORT := os.Getenv("PORT")
		if PORT == "" { PORT = "80" }
		Config = IConfig{
			"PORT": PORT,
		}
	}
	return Config
}