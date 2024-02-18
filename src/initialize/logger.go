package initialize

import (
	"app/config"
	"app/di"
	"app/utils"
	"fmt"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func InitLogger(){
	utils.MakeSureDir("tmp")
	host, _ := os.Hostname()
	if host == "" {
		host = "app"
	}

	var cfg = zap.Config{
		Encoding: "json",
		OutputPaths: []string{"stdout", path.Join("tmp", fmt.Sprintf("%s-%s.log", config.Config["GO_ENV"], host))},
		ErrorOutputPaths: []string{"stderr", path.Join("tmp", fmt.Sprintf("error-%s-%s.log", config.Config["GO_ENV"], host))},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	switch config.Config["GO_ENV"] {
	case "development","test":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger := zap.Must(cfg.Build())
	di.Container.Logger = logger
}