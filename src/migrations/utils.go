package migrations

import (
	"app/di"
	"app/initialize"
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

func getDB(ctx context.Context, tx *sql.Tx) *gorm.DB {
	if di.Container.DB != nil {
		return di.Container.DB
	}
	initialize.InitDB()
	if di.Container.DB == nil {
		panic(fmt.Errorf("connect db failed"))
	}
	return di.Container.DB
}

func createTable(migrator gorm.Migrator, dst interface{}) {
	if migrator.HasTable(dst) {
		return
	}
	if err := migrator.CreateTable(dst); err != nil {
		panic(err)
	}
}
func dropTable(migrator gorm.Migrator, dst interface{}) {
	if err := migrator.DropTable(dst); err != nil {
		panic(err)
	}
}

func addTableColumn(migrator gorm.Migrator, dst interface{}, field string) {
	if migrator.HasColumn(dst, field) {
		return
	}
	if err := migrator.AddColumn(dst, field); err != nil {
		panic(err)
	}
}
func dropTableColumn(migrator gorm.Migrator, dst interface{}, field string) {
	if !migrator.HasColumn(dst, field) {
		return
	}
	if err := migrator.DropColumn(dst, field); err != nil {
		panic(err)
	}
}
func alertTableColumn(migrator gorm.Migrator, dst interface{}, field string) {
	if !migrator.HasColumn(dst, field) {
		return
	}
	if err := migrator.AlterColumn(dst, field); err != nil {
		panic(err)
	}
}

func createTableIndex(migrator gorm.Migrator, dst interface{}, field string) {
	if migrator.HasIndex(dst, field) {
		return
	}
	if err := migrator.CreateIndex(dst, field); err != nil {
		panic(err)
	}
}
func dropTableIndex(migrator gorm.Migrator, dst interface{}, field string) {
	if !migrator.HasIndex(dst, field) {
		return
	}
	if err := migrator.DropIndex(dst, field); err != nil {
		panic(err)
	}
}