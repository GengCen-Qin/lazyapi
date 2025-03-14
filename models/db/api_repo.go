package db

import (
	"database/sql"
	"lazyapi/models/entity"
)

func InsertAPI(api *entity.API) error {
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

func UpdateAPI(api *entity.API) error {
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

func FindAPI(id int) (*entity.API, error) {
	db := GetDB()
    stmt, err := db.Prepare("SELECT id, name, path, method, params FROM apis WHERE id=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    var api entity.API
    err = stmt.QueryRow(id).Scan(&api.Id, &api.Name, &api.Path, &api.Method, &api.Params)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // 没有找到记录，返回nil而不是错误
        }
        return nil, err
    }

    return &api, nil
}

func GetAllAPIs() ([]entity.API, error) {
	db := GetDB()
    rows, err := db.Query("SELECT id, name, path, method, params FROM apis")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var apis []entity.API
    for rows.Next() {
        var api entity.API
        if err := rows.Scan(&api.Id, &api.Name, &api.Path, &api.Method, &api.Params); err != nil {
            return nil, err
        }
        apis = append(apis, api)
    }
    return apis, nil
}
