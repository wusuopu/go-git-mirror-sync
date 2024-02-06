package migrations

import (
	"app/utils"
	"context"
	"database/sql"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// 20240206071430-CreateRepository

func init() {
	goose.AddMigrationContext(up20240206071430, down20240206071430)
}
func createModel20240206071430 () interface{} {
	type Repository struct{
		gorm.Model
		Name        string      `gorm:"type:varchar(80);"`
		Alias       string      `gorm:"type:varchar(80);"`
		Url         string      `gorm:"type:varchar(255);"`
		AuthType    string      `gorm:"type:varchar(15);"`      // password | sshkey
		Username    *string     `gorm:"type:varchar(80);"`
		Password    *string     `gorm:"type:varchar(80);"`
		SSHKey      *string     `gorm:"type:varchar(80);"`
		PulledAt    *time.Time
	}
	return &Repository{}
}
func up20240206071430(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20240206071430()
		createTable(migrator, model)
	})
}

func down20240206071430(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20240206071430()
		dropTable(migrator, model)
	})
}
