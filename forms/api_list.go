package forms

import (
	"fmt"
	"slices"

	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/models/service"

	"github.com/GengCen-Qin/gocui"
)

// UpdateAPIList 更新左侧API列表显示
func UpdateAPIList(g *gocui.Gui) {
	leftView, leftViewError := g.View("left")
	if leftViewError != nil {
		return // 如果视图不存在，直接返回
	}
	rightTopView, _ := g.View("right-top")
	list := service.APIList()
	leftView.Clear()
	for _, api := range list {
		if api.Id == entity.SelectedAPI {
			fmt.Fprintf(leftView, "> %s [%s] \n", api.Name, api.Method)
		} else {
			// 文字颜色控制
			fmt.Fprintf(leftView, "  %s [\033[31;1m%s\033[0m] \n", api.Name, api.Method)
		}
	}

	index := slices.IndexFunc(list, func(x entity.API) bool {
	    return x.Id == entity.SelectedAPI
	})

	rightTopView.Clear()

	if entity.SelectedAPI != -1 {
		leftView.SetCursor(0, index)
		api, _ := db.FindAPI(entity.SelectedAPI)
		fmt.Fprintf(rightTopView, "\033[34;1mName\033[0m: %s \t \033[34;1mMethod\033[0m: %s\n",api.Name, api.Method)
		fmt.Fprintf(rightTopView, "\033[34;1mPath\033[0m: %s\n", api.Path)
		fmt.Fprintf(rightTopView, "\033[34;1mParams\033[0m: \n%s\n", api.Params)
	} else {
		fmt.Fprint(rightTopView, "EMPTY API")
	}
}

// RefreshRequestRecordList 刷新请求记录列表
func RefreshRequestRecordList(g *gocui.Gui) {
	view, _ := g.View("request-history")

	list := service.RequestRecordList()

	view.Clear()
	for _, record := range list {
		if record.Id == entity.SelectedQuestRecord {
		    fmt.Fprintf(view, ">\033[34;1m%s\033[0m [\033[a31;1m%s\033[0m] \n",
		                record.RequestTime.Local().Format("2006-01-02 15:04:05"), record.Path)
		} else {
		    fmt.Fprintf(view, " \033[34;1m%s\033[0m [\033[31;1m%s\033[0m] \n",
		                record.RequestTime.Local().Format("2006-01-02 15:04:05"), record.Path)
		}
	}
}

// MoveSelectionUp 向上移动选择
func MoveSelectionUp(g *gocui.Gui, v *gocui.View) error {
	list := service.APIList()
	if len(list) == 0 || len(list) == 1 || entity.SelectedAPI <= 0 {
		return nil
	}

	index := slices.IndexFunc(list, func(x entity.API) bool {
	    return x.Id == entity.SelectedAPI
	})

	if index != 0 {
		entity.SelectedAPI = list[index-1].Id
	}
	UpdateAPIList(g)
	return nil
}

// MoveSelectionDown 向下移动选择
func MoveSelectionDown(g *gocui.Gui, v *gocui.View) error {
	list := service.APIList()
	if len(list) == 0 || len(list) == 1 {
		return nil
	}

	index := slices.IndexFunc(list, func(x entity.API) bool {
	    return x.Id == entity.SelectedAPI
	})

	if index + 1 >= len(list) {
		return nil
	}

	entity.SelectedAPI = list[index+1].Id

	UpdateAPIList(g)
	return nil
}
