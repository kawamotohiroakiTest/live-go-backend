package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// データベース接続情報とSQLファイルのパスを定義
const (
	Source   = "file://db/migrations"
	Database = "mysql://liveuser:livepass@tcp(db:3306)/live"
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

func main() {
	RunMigration()
}

func RunMigration() {
	flag.Parse()

	// command引数が指定されていない場合、デフォルトで "up" コマンドを実行
	if len(*Command) < 1 {
		*Command = "up"
		fmt.Println("No command provided, defaulting to 'up' migration.")
	}

	m, err := migrate.New(Source, Database)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	// 現在のバージョン情報を取得
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get current migration version: %v", err)
	}

	// "up" コマンドの場合、更新がなければマイグレーションを実行しない
	if *Command == "up" {
		latestVersion := getLatestVersion() // 最新のマイグレーションファイルのバージョンを取得
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
	files, err := ioutil.ReadDir("./db/migrations")
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var latestVersion uint

	for _, file := range files {
		// ファイル名が .up.sql で終わるものを対象にする
		if strings.HasSuffix(file.Name(), ".up.sql") {
			// ファイル名の先頭部分からバージョン番号を取得する
			versionStr := strings.Split(file.Name(), "_")[0]
			version, err := strconv.ParseUint(versionStr, 10, 64)
			if err != nil {
				log.Printf("Failed to parse migration version from file %s: %v", file.Name(), err)
				continue
			}

			// 最新のバージョン番号を更新
			if uint(version) > latestVersion {
				latestVersion = uint(version)
			}
		}
	}

	if latestVersion == 0 {
		log.Fatalf("No valid migration files found in the directory")
	}

	fmt.Printf("Latest migration version: %d\n", latestVersion)
	return latestVersion
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
		// do nothing
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
		fmt.Println("Updated version info")
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
