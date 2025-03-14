package db

import (
    "database/sql"
    "log"
    "os"
    "path/filepath"
    "runtime"

    _ "github.com/mattn/go-sqlite3"
)

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
        log.Println("数据库连接已关闭")
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
    query := `
    CREATE TABLE IF NOT EXISTS apis (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        path TEXT NOT NULL,
        method TEXT NOT NULL,
        params TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS request_records (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        path TEXT NOT NULL,
        method TEXT NOT NULL,
        params TEXT NOT NULL,
        respond TEXT NOT NULL,
        request_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
    _, err := db.Exec(query)
    if err != nil {
        log.Println("Error initializing database:", err)
    } else {
        log.Println("Database initialized successfully.")
    }
    return err
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
