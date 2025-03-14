package entity

import "time"

type RequestRecord struct {
    Id          int       `json:"id"`
    Name        string    `json:"name"`
    Path        string    `json:"path"`
    Method      string    `json:"method"`
    Params      string    `json:"params"`
    Respond     string    `json:"respond"`
    RequestTime time.Time `json:"request_time"`
}

var SelectedQuestRecord int = -1
