package config

import (
	"os"
	"strconv"
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
		GIT_INSECURE_SKIP_TLS, _ := strconv.ParseBool(os.Getenv("GIT_INSECURE_SKIP_TLS"))

		CRONTAB := os.Getenv("CRONTAB")
		if CRONTAB == "" {
			// 每6小时执行一次
			CRONTAB = "0 */6 * * *"
		}

		Config = IConfig{
			"PORT": PORT,
			"GO_ENV": GO_ENV,
			"GIT_INSECURE_SKIP_TLS": GIT_INSECURE_SKIP_TLS,
			"CRONTAB": CRONTAB,
		}
	}
	return Config
}
