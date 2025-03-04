package models

// API 表示一个API接口的数据结构
type API struct {
	Name   string `json:"name"`   // API名称
	Path   string `json:"path"`   // API路径
	Method string `json:"method"` // HTTP方法
}

// NewAPI 创建一个新API
func NewAPI(name, path, method string) *API {
	return &API{
		Name:   name,
		Path:   path,
		Method: method,
	}
}

// APIList 存储所有API
var APIList []*API
var SelectedAPI int = -1
