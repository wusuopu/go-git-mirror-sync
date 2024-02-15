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
// 20240207042557-AddInitedatOnRepository

func init() {
	goose.AddMigrationContext(up20240207042557, down20240207042557)
}
func createModel20240207042557 () interface{} {
	type Repository struct{
		gorm.Model
		InitedAt    *time.Time
		LastError   *string     `gorm:"type:text;"`
	}
	return &Repository{}
}
func up20240207042557(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20240207042557()
		addTableColumn(migrator, model, "InitedAt")
		addTableColumn(migrator, model, "LastError")
	})
}

func down20240207042557(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20240207042557()
		dropTableColumn(migrator, model, "InitedAt")
		dropTableColumn(migrator, model, "LastError")
	})
}
