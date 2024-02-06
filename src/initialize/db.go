package initialize

import (
	"app/di"
	"database/sql"
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
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

// 根据已有的 sql 链接创建 GORM 对象
func DBFromConn(conn *sql.DB, driver string) (*gorm.DB, error) {
	currentDriver := driver
	if currentDriver == "" {
		currentDriver = os.Getenv("DATABASE_TYPE")
	}
	if currentDriver == "mysql" {
		dia := mysql.New(mysql.Config{
			Conn: conn,
		})
		return gorm.Open(dia, &gorm.Config{})
	}

	dia := sqlite.Dialector{Conn: conn}
	return gorm.Open(&dia, &gorm.Config{})
}

func ConnectMySQLWithoutDB(dsn string) (*gorm.DB, string, error) {
	dialector := mysql.Open(dsn)
	dbName := dialector.(*mysql.Dialector).DSNConfig.DBName
	dialector.(*mysql.Dialector).DSNConfig.DBName = ""
	newDSN := dialector.(*mysql.Dialector).DSNConfig.FormatDSN()

	dialector = mysql.New(mysql.Config{
		DSN: newDSN,
		DefaultStringSize:         255,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置

	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	return db, dbName, err
}