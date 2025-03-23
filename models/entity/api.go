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

var (
	SelectedAPI int = -1

	Method_Unkown      = 0
	Method_Get         = 1
	Method_Post        = 2
	MethodTitle = map[int]string{
		0: "未知",
		1: "GET",
		2: "POST",
	}
)
