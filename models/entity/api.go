package entity

import "encoding/json"

type API struct {
    Id     int    `json:"id"`
    Name   string `json:"name"`
    Path   string `json:"path"`
    Method string `json:"method"`
    Params string `json:"params"`
}

func (api *API) GetParams() (map[string]interface{}, error) {
    var params map[string]interface{}
    err := json.Unmarshal([]byte(api.Params), &params)
    if err != nil {
        return nil, err
    }
    return params, nil
}

var SelectedAPI int = -1
