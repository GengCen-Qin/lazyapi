package db

import (
	"lazyapi/models/entity"
)

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
    rows, err := db.Query("SELECT id, name, path, method, params, respond, request_time FROM request_records")
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
