package notes

import (
	"app/config"
	"app/di"
	"app/initialize"
	"app/utils"
	"fmt"
	"os"
	"path"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	var dialector gorm.Dialector
	if os.Getenv("DATABASE_TYPE") == "mysql" {
		dialector = mysql.New(mysql.Config{
			DSN: os.Getenv("DATABASE_DSN"),
			DefaultStringSize:         255,     // string 类型字段的默认长度
			SkipInitializeWithVersion: false,   // 根据版本自动配置

		})
	} else {
		// https://gorm.io/docs/connecting_to_the_database.html#SQLite
		dialector = sqlite.Open(os.Getenv("DATABASE_DSN"))
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("connect DB error: %s", err.Error()))
	}
	di.Container.DB = db
}
func InitLogger(){
	utils.MakeSureDir("tmp")

	var cfg = zap.Config{
		Encoding: "json",
		OutputPaths: []string{"stdout", path.Join("tmp", fmt.Sprintf("%s.log", config.Config["GO_ENV"]))},
		ErrorOutputPaths: []string{"stderr", path.Join("tmp", fmt.Sprintf("error-%s.log", config.Config["GO_ENV"]))},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
			LevelKey: "level",
		},
	}

	switch config.Config["GO_ENV"] {
	case "development","test":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger := zap.Must(cfg.Build())
	di.Container.Logger = logger
}



func Init(filename string) {
	godotenv.Load(filename)
	config.Load()
	InitDB()
	InitLogger()
	initialize.InitServices()
}