package forms

import (
	"fmt"
	"strings"

	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/models/service"

	"github.com/GengCen-Qin/gocui"
)

// SaveNewAPI 保存API（新建或编辑）
func SaveNewAPI(g *gocui.Gui, v *gocui.View) error {
	if !FormInfo.Active {
		return nil
	}
	if FormInfo.IsDelete {
		return nil
	}

	// 收集表单数据
	var name, path, method, params string
	nameView, _ := g.View("form-name")
	pathView, _ := g.View("form-path")
	methodView, _ := g.View("form-method")
	paramsView, _ := g.View("form-params")

	name = strings.TrimSpace(nameView.Buffer())
	path = strings.TrimSpace(pathView.Buffer())
	method = strings.TrimSpace(methodView.Buffer())
	params = strings.TrimSpace(paramsView.Buffer())

	if err := validateAPIForm(g, name, path, method, params); err != nil {
        statusView, _ := g.View("status")
        statusView.Clear()
        fmt.Fprint(statusView, err.Error())
        return nil
    }

	// 如果是编辑模式，更新现有API
	if FormInfo.IsEditing {
		service.EditAPI(entity.SelectedAPI, name, path, method, params)
	} else {
		// 否则，创建新API并添加到列表
		newAPI := service.NewAPI(name, path, method, params)
		entity.SelectedAPI = newAPI.Id
	}

	// 更新视图
	UpdateAPIList(g)

	// 关闭表单
	FormInfo.IsEditing = false
	// 关闭光标
	g.Cursor = false
	return CloseForm(g, v)
}

// EditAPIForm 编辑API表单
func EditAPIForm(g *gocui.Gui, v *gocui.View) error {
	if entity.SelectedAPI == -1 {
		return nil
	}

	// 标记为编辑模式
	FormInfo.IsEditing = true

	// 显示表单并填充数据
	if err := ShowNewAPIForm(g, v); err != nil {
		return err
	}

	// 获取当前选中的API
	api, _ := db.FindAPI(entity.SelectedAPI)

	fillFormFields(g, api)
	return nil
}
