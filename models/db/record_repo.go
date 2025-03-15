package db

import (
	"database/sql"
	"lazyapi/models/entity"
)

func Find(id int) (*entity.RequestRecord, error) {
	db := GetDB()
    stmt, err := db.Prepare("SELECT id, name, path, method, params, respond, request_time FROM request_records WHERE id=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    var record entity.RequestRecord
    err = stmt.QueryRow(id).Scan(&record.Id, &record.Name, &record.Path, &record.Method, &record.Params, &record.Respond, &record.RequestTime)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // 没有找到记录，返回nil而不是错误
        }
        return nil, err
    }

    return &record, nil
}

func DeleteRecord(id int) error {
	db := GetDB()
    stmt, err := db.Prepare("DELETE FROM request_records WHERE id=?")
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

func InsertRequestRecord(api *entity.API, params string, respond string) error {
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


func GetAllRequestRecords() ([]entity.RequestRecord, error) {
    db := GetDB()
    rows, err := db.Query("SELECT id, name, path, method, params, respond, request_time FROM request_records order by request_time desc")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []entity.RequestRecord
    for rows.Next() {

        var request_record entity.RequestRecord
        if err := rows.Scan(&request_record.Id, &request_record.Name, &request_record.Path, &request_record.Method, &request_record.Params, &request_record.Respond, &request_record.RequestTime); err != nil {
            return nil, err
        }
        list = append(list, request_record)
    }
    return list, nil
}
