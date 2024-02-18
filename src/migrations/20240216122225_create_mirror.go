package migrations

import (
	"app/models"
	"app/utils"
	"context"
	"database/sql"
	"time"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// 20240216122225-CreateMirror

func init() {
	goose.AddMigrationContext(up20240216122225, down20240216122225)
}
func createModel20240216122225 () interface{} {
	type Mirror struct{
		gorm.Model
		Name        string      `gorm:"type:varchar(80);"`
		Alias       string      `gorm:"type:varchar(80);"`
		Url         string      `gorm:"type:varchar(255);"`
		AuthType    string      `gorm:"type:varchar(15);"`      // password | sshkey
		Username    *string     `gorm:"type:varchar(80);"`
		Password    *string     `gorm:"type:varchar(80);"`
		SSHKey      *string     `gorm:"type:text;"`
		PushedAt    *time.Time
		LastError   *string     `gorm:"type:text;"`
		RepositoryId uint
		Repository  models.Repository
	}
	return &Mirror{}
}
func up20240216122225(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20240216122225()
		createTable(migrator, model)
	})
}

func down20240216122225(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20240216122225()
		dropTable(migrator, model)
	})
}
