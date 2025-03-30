package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type silentLogger struct{}

func (s *silentLogger) Fatalf(format string, v ...interface{}) {
}

func (s *silentLogger) Printf(format string, v ...interface{}) {
}

var dbInstance *sql.DB

func GetDB() *sql.DB {
    if dbInstance == nil {
        initDB()
    }
    return dbInstance
}

func CloseDB() {
    if dbInstance != nil {
        dbInstance.Close()
        dbInstance = nil
    }
}

func initDB() {
    var err error
    dbInstance, err = sql.Open("sqlite3", dbPath())
    if err != nil {
        log.Fatal(err)
    }

    if err = dbInstance.Ping(); err != nil {
        log.Fatal(err)
    }

    createTables()
}

func createTables() error {
    db := GetDB()

    workDir, _ := os.Getwd()
    migrationsDir := filepath.Join(workDir, "migrations")

    // 设置 goose 使用静默的日志记录器
    goose.SetLogger(&silentLogger{})
    goose.SetDialect("sqlite3")

    err := goose.Up(db, migrationsDir)
    if err != nil {
        return err
    }

    return nil
}

func dbPath() string {
	var dbPath string
	homeDir, err := os.UserHomeDir()
	if err != nil {
	    // 如果无法获取用户主目录，则使用当前目录作为备选
	    dbPath = "./lazyapi.db"
	} else {
	    if runtime.GOOS == "darwin" { // Mac系统
	        // 创建Mac应用支持目录
	        appSupportDir := filepath.Join(homeDir, "Library", "Application Support", "lazyapi")
	        if err := os.MkdirAll(appSupportDir, 0755); err != nil {
	            // 如果无法创建目录，回退到用户主目录
	            dbPath = filepath.Join(homeDir, "lazyapi.db")
	        } else {
	            dbPath = filepath.Join(appSupportDir, "lazyapi.db")
	        }
	    } else {
	        // 其他系统默认放在用户主目录下
	        dbPath = filepath.Join(homeDir, "lazyapi.db")
	    }
	}
	return dbPath
}
