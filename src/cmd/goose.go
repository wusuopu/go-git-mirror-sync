package main

import (
	"app/initialize"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"text/template"

	_ "app/migrations"

	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	env = flags.String("env", "development", "当前的环境： development | production | test")
)
func confirm(msg string) bool {
	fmt.Printf("%s  [y/N]?", msg)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return strings.Trim(strings.ToLower(text), "\n") == "y"
}
func commandDBCreate(driver, dsn string) bool {
	if driver == "sqlite" {
		return true
	}
	db, dbName, err := initialize.ConnectMySQLWithoutDB(dsn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", dbName))
	fmt.Printf("数据库 %s 创建成功\n", dbName)
	return true
}
func commandDBDrop(driver, dsn string) bool {
	if !confirm(fmt.Sprintf("是否删除数据库 %s %s", driver, dsn)) {
		fmt.Println("数据库未删除")
		return false
	}
	if driver == "sqlite" {
		u, err := url.Parse(dsn)
		if err != nil {
			return false
		}
		info, err := os.Stat(u.Path)
		if err != nil {
			return false
		}
		if info.IsDir() {
			return false
		}
		err = os.Remove(u.Path)
		if err != nil {
			fmt.Println(err)
			return false
		}
		fmt.Printf("%s 文件删除成功\n", u.Path)
		return true
	}

	db, dbName, err := initialize.ConnectMySQLWithoutDB(dsn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", dbName))
	fmt.Printf("数据库 %s 删除成功\n", dbName)
	return true
}
func commandMigrationCreate(args ...string) bool {
	if len(args) == 0 {
		fmt.Printf("create must be of form: goose [OPTIONS] create NAME\n")
		return false
	}

	tpl := template.Must(template.New("goose.go-migration").Parse(`package migrations

import (
	"app/utils"
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// https://gorm.io/docs/migration.html
// {{.Version}}-{{.CamelName}}

func init() {
	goose.AddMigrationContext(up{{.Version}}, down{{.Version}})
}
func createModel{{.Version}} () interface{} {
	// TODO 修改对应的 ModelName
	type ModelName struct{
		gorm.Model
	}
	return &ModelName{}
}
func up{{.Version}}(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is applied.
		migrator := db.Migrator()
		// TODO 根据情况修改
		model := createModel{{.Version}}()
		addTableColumn(migrator, model, "FieldName")
	})
}

func down{{.Version}}(ctx context.Context, tx *sql.Tx) error {
	return utils.Try(func() {
		db := getDB(ctx, tx)

		// This code is executed when the migration is rolled back.
		migrator := db.Migrator()
		model := createModel{{.Version}}()
		// TODO 根据情况修改
		dropTableColumn(migrator, model, "FieldName")
	})
}
`))

	goose.CreateWithTemplate(
		nil,
		"migrations",
		tpl,
		args[0],
		"go",
	)
	return true
}


var usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME          Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
    validate             Check migration files without running them
    db:create            create database
    db:drop              drop database
`
func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if *env == "test" {
		initialize.InitEnv(".env.test")
	} else {
		initialize.InitEnv()
	}

	if len(args) < 1 {
		flags.Usage()
		fmt.Println("goose command [options]")
		fmt.Println(usageCommands)
		return
	}
	driver := os.Getenv("DATABASE_TYPE")
	if driver != "mysql" {
		driver = "sqlite"
	}
	fmt.Printf("当前环境: %s, 数据库: %s\n", *env, driver)

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	// goose.SetBaseFS(embedMigrations)
	command := args[0]
	switch command {
	case "db:create":
		commandDBCreate(driver, os.Getenv("DATABASE_DSN"))
		return
	case "db:drop":
		commandDBDrop(driver, os.Getenv("DATABASE_DSN"))
		return
	case "create":
		commandMigrationCreate(arguments...)
		return
	}


	db, err := goose.OpenDBWithDriver(driver, os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func ()  {
		db.Close()	
	}()
	variables := map[string]interface{}{
		"driver": driver,
	}
	ctx := context.WithValue(context.Background(), "args", variables)
	// goose.SetTableName("goose_db_version")		// 默认使用 goose_db_version 表记录操作
	err = goose.RunContext(ctx, command, db, ".", arguments...)
	if err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}