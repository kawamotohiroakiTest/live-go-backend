package db

import (
	"flag"
	"fmt"
	"io/ioutil"
	"live/common"
	"log"
	"os"
	"strconv"
	"strings"

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

	source := os.Getenv("MIGRATION_SOURCE")
	if source == "" {
		source = "file://db/migrations"
	}

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	fmt.Printf("Connecting to MySQL with DSN: %s\n", dsn)

	if *Command == "" {
		*Command = "up" // デフォルトコマンドを設定
	}

	m, err := migrate.New(source, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get current migration version: %v", err)
	}

	// getLatestVersionから2つの値を受け取る
	latestVersion, latestFile := getLatestVersion()

	if version == latestVersion && !dirty {
		fmt.Println("No new migrations to apply.")
		fmt.Printf("Current migration version: %d, Latest migration version: %d\n", version, latestVersion)
		return
	}
	// 最新のマイグレーションファイルとそのバージョンを表示
	fmt.Printf("Latest migration version: %d, File: %s\n", latestVersion, latestFile)

	applyQuery(m, version, dirty)
}

func getLatestVersion() (uint, string) {
	// db/migrationsディレクトリのファイルを読み込む
	files, err := ioutil.ReadDir("db/migrations")
	if err != nil {
		log.Fatalf("Failed to read migration directory: %v", err)
	}

	var latestVersion uint
	var latestFile string

	for _, file := range files {
		// ファイル名からバージョン番号を取得（数値部分を抽出）
		filename := file.Name()
		versionStr := strings.Split(filename, "_")[0]
		version, err := strconv.ParseUint(versionStr, 10, 64)
		if err != nil {
			log.Printf("Skipping file %s: %v", filename, err)
			continue
		}

		// 最も新しいバージョンを更新
		if uint(version) > latestVersion {
			latestVersion = uint(version)
			latestFile = filename
		}
	}

	// 最新のバージョン番号とファイル名を返す
	return latestVersion, latestFile
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
