package models

import (
	"encoding/json"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// API 表示一个API接口的数据结构
type API struct {
	Id     int    `json:"id"`     // API ID
	Name   string `json:"name"`   // API名称
	Path   string `json:"path"`   // API路径
	Method string `json:"method"` // HTTP方法
	Params string `json:"params"` // 请求参数，存储为JSON字符串
}

// NewAPI 创建一个新API
func NewAPI(name, path, method string, params string) *API {
	params = formatJSON(params)
    
    api := &API{
		Name:   name,
		Path:   path,
		Method: method,
		Params: params,
	}

    InsertAPI(api)

	return api
}

func EditAPI(id int, name, path, method string, params string) *API {
	params = formatJSON(params)

	api, _ := FindAPI(id)
	api.Name = name
	api.Path = path
	api.Method = method
	api.Params = params
	UpdateAPI(api)

	return api
}

func formatJSON(jsonString string) string {
    var jsonObj map[string]interface{}

    // Parse the JSON into a map
    err := json.Unmarshal([]byte(jsonString), &jsonObj)
    if err != nil {
        log.Fatalf("Error occured during unmarshalling. %s", err)
    }

    // Format the json
    formattedJSON, err := json.MarshalIndent(jsonObj, "", "    ")
    if err != nil {
        log.Fatalf("Error occured during marshalling. %s", err)
    }

    return string(formattedJSON)
}

// GetParams 获取请求参数
func (a *API) GetParams() (map[string]interface{}, error) {
	var params map[string]interface{}
	err := json.Unmarshal([]byte(a.Params), &params)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func APIList() ([]API) {
    apis, _ := getAllAPIs()
	return apis
}

var SelectedAPI int = -1
