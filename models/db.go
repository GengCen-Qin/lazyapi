package models

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB // 全局数据库连接池

func init() {
	dbPath := "./lazyapi.db"
    var err error
    db, err = sql.Open("sqlite3", dbPath)
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