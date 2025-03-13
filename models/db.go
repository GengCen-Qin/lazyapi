package models

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB // 全局数据库连接池

func init() {
    var err error
    db, err = sql.Open("sqlite3", dbPath())
    if err != nil {
        log.Fatal(err)
    }

    // 测试连接是否有效
    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    err = initDB(db)
    if err != nil {
        log.Fatal(err)
    }
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

func CloseDB() {
    if db != nil {
        db.Close()
        log.Println("数据库连接已关闭")
    }
}

func GetDB() *sql.DB {
    return db
}

func initDB(db *sql.DB) error {
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

func InsertAPI(api *API) error {
	db := GetDB()
    stmt, err := db.Prepare("INSERT INTO apis(name, path, method, params) values(?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    result, err := stmt.Exec(api.Name, api.Path, api.Method, api.Params)
    if err != nil {
        return err
    }
    lastID, err := result.LastInsertId()
    if err != nil {
        return err
    }
    api.Id = int(lastID)
    return nil
}

func InsertRequestRecord(api *API, params string, respond string) error {
    db := GetDB()
    stmt, err := db.Prepare("INSERT INTO request_records(name, path, method, params, respond) values(?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    _, err = stmt.Exec(api.Name, api.Path, api.Method, params, respond)
    if err != nil {
        return err
    }

    return nil
}

func UpdateAPI(api *API) error {
	db := GetDB()
    stmt, err := db.Prepare("UPDATE apis SET name=?, path=?, method=?, params=? WHERE id=?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(api.Name, api.Path, api.Method, api.Params, api.Id)
    if err != nil {
        return err
    }

    return nil
}

func DeleteAPI(id int) error {
	db := GetDB()
    stmt, err := db.Prepare("DELETE FROM apis WHERE id=?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        return err
    }

    return nil
}

func FindAPI(id int) (*API, error) {
	db := GetDB()
    stmt, err := db.Prepare("SELECT id, name, path, method, params FROM apis WHERE id=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    var api API
    err = stmt.QueryRow(id).Scan(&api.Id, &api.Name, &api.Path, &api.Method, &api.Params)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // 没有找到记录，返回nil而不是错误
        }
        return nil, err
    }

    return &api, nil
}

func getAllAPIs() ([]API, error) {
	db := GetDB()
    rows, err := db.Query("SELECT id, name, path, method, params FROM apis")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var apis []API
    for rows.Next() {
        var api API
        if err := rows.Scan(&api.Id, &api.Name, &api.Path, &api.Method, &api.Params); err != nil {
            return nil, err
        }
        apis = append(apis, api)
    }
    return apis, nil
}

func getAllRequestRecords() ([]RequestRecord, error) {
    db := GetDB()
    rows, err := db.Query("SELECT id, name, path, method, params, respond, request_time FROM request_records")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []RequestRecord
    for rows.Next() {

        var request_record RequestRecord
        if err := rows.Scan(&request_record.Id, &request_record.Name, &request_record.Path, &request_record.Method, &request_record.Params, &request_record.Respond, &request_record.RequestTime); err != nil {
            return nil, err
        }
        list = append(list, request_record)
    }
    return list, nil
}
