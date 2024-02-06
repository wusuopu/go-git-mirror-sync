package notes

import (
	"app/config"
	"app/di"
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
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
func Init(filename string) {
	godotenv.Load(filename)
	config.Load()
	InitDB()
}