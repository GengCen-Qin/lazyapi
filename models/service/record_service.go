package service

import (
	"lazyapi/models/entity"
	"lazyapi/models/db"
)

func RequestRecordList() ([]entity.RequestRecord) {
	records, _ := db.GetAllRequestRecords()
	return records
}
