package service

import (
	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/utils"
)

func NewAPI(name, path, method string, params string) *entity.API {
	params = utils.FormatJSON(params)

    api := &entity.API{
		Name:   name,
		Path:   path,
		Method: method,
		Params: params,
	}

    db.InsertAPI(api)

	return api
}

func EditAPI(id int, name, path, method string, params string) *entity.API {
	params = utils.FormatJSON(params)

	api, _ := db.FindAPI(id)
	api.Name = name
	api.Path = path
	api.Method = method
	api.Params = params
	db.UpdateAPI(api)

	return api
}

func APIList() ([]entity.API) {
    apis, _ := db.GetAllAPIs()
	return apis
}
