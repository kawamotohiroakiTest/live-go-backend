package db

import (
	"flag"
	"fmt"
	"live/common"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// コマンドラインオプションの宣言
var (
	Command = flag.String("exec", "", "Specify up, down, or version as an argument")
	Force   = flag.Bool("f", false, "Force execute the migration")
)

// 実行可能なコマンドのリスト
var AvailableExecCommands = map[string]string{
	"up":      "Execute up migrations",
	"down":    "Execute down migrations",
	"version": "Check current migration version",
}

func RunMigration() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		common.LogError(fmt.Errorf("Error loading .env file: %v", err))
	}

	// 環境変数からSourceとDatabaseの情報を取得
	source := os.Getenv("MIGRATION_SOURCE")
	if source == "" {
		source = "file://db/migrations"
	}

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	// MySQLコンテナへの接続情報
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	// DSNの確認
	fmt.Printf("Connecting to MySQL with DSN: %s\n", dsn)

	// command引数が指定されていない場合、デフォルトで "up" コマンドを実行
	if len(*Command) < 1 {
		*Command = "up"
		fmt.Println("No command provided, defaulting to 'up' migration.")
	}

	m, err := migrate.New(source, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	// マイグレーションの実行などの処理
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get current migration version: %v", err)
	}

	if *Command == "up" {
		latestVersion := getLatestVersion()
		if version == latestVersion && !dirty {
			fmt.Println("No new migrations to apply.")
			return
		}
	}

	fmt.Println("Command: exec", *Command)
	applyQuery(m, version, dirty)
}

// 最新のマイグレーションファイルのバージョンを取得するヘルパー関数
func getLatestVersion() uint {
	// 実際のロジックでは、マイグレーションディレクトリから最新のファイルのバージョンを取得するようにします。
	return 20240826142057
}

// マイグレーションを実行する関数
func applyQuery(m *migrate.Migrate, version uint, dirty bool) {
	if dirty && *Force {
		fmt.Println("Force=true: Force execute current version migration")
		if err := m.Force(int(version)); err != nil {
			log.Fatalf("Failed to force migration: %v", err)
		}
	}

	var err error
	switch *Command {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "version":
		return
	default:
		fmt.Println("\nError: Invalid command '" + *Command + "'\n")
		showUsageMessage()
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		fmt.Printf("Migration error: %v\n", err)
		os.Exit(1)
	} else if err == migrate.ErrNoChange {
		fmt.Println("No new migrations to apply.")
	} else {
		fmt.Println("Success:", *Command)
		version, dirty, err := m.Version()
		showVersionInfo(version, dirty, err)
	}
}

// 使用方法を表示する関数
func showUsageMessage() {
	fmt.Println("-------------------------------------")
	fmt.Println("Usage:")
	fmt.Println("  go run main.go -exec <command>")
	fmt.Println("\nAvailable Exec Commands:")
	for availableCommand, detail := range AvailableExecCommands {
		fmt.Println("  " + availableCommand + " : " + detail)
	}
	fmt.Println("-------------------------------------")
}

// マイグレーションのバージョン情報を表示する関数
func showVersionInfo(version uint, dirty bool, err error) {
	fmt.Println("-------------------")
	fmt.Printf("Version  : %d\n", version)
	fmt.Printf("Dirty    : %v\n", dirty)
	fmt.Printf("Error    : %v\n", err)
	fmt.Println("-------------------")
}
