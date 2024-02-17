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
// 20240216123200-CreateBranch

func init() {
	goose.AddMigrationContext(up20240216123200, down20240216123200)
}
func createModel20240216123200 () interface{} {
	type Branch struct {
		gorm.Model
		Name        string      `gorm:"type:varchar(80);"`
		Hash        string      `gorm:"type:varchar(80);"`
		CommittedAt *time.Time
		CommitMsg   *string     `gorm:"type:text;"`
		IsTag       bool
		RepositoryId	uint
		Repository	models.Repository
	}
	return &Branch{}
}
func up20240216123200(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		model := createModel20240216123200()
		createTable(migrator, model)
	})
}

func down20240216123200(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel20240216123200()
		dropTable(migrator, model)
	})
}
