package forms

import (
	"bytes"
	"fmt"
	"slices"

	"lazyapi/common"
	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/models/service"
	"lazyapi/utils"

	"github.com/GengCen-Qin/gocui"
)

// UpdateAPIList 更新左侧API列表显示
func UpdateAPIList(g *gocui.Gui) {
	leftView, leftViewError := g.View("left")
	if leftViewError != nil {
		return // 如果视图不存在，直接返回
	}
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

	if common.ViewActiveIndex != common.ViewIndexMap["left"] {
    	return
    }

	index := slices.IndexFunc(list, func(x entity.API) bool {
	    return x.Id == entity.SelectedAPI
	})

	rightTopView, _ := g.View("right-top")
	rightTopView.Clear()
	rightBottomView, _ := g.View("right-bottom")
	rightBottomView.Clear()

	if entity.SelectedAPI != -1 {
		leftView.SetCursor(0, index)
		api, _ := db.FindAPI(entity.SelectedAPI)
		var buffer bytes.Buffer
	    format_json, _ := utils.PrettyPrintJSON(api.Params)
		fmt.Fprintf(&buffer, "\033[34;1mName\033[0m: %s \t \033[34;1mMethod\033[0m: %s\n", api.Name, api.Method)
		fmt.Fprintf(&buffer, "\033[34;1mPath\033[0m: %s\n", api.Path)
		fmt.Fprintf(&buffer, "\033[34;1mParams\033[0m: \n%s\n", format_json)
		fmt.Fprint(rightTopView, buffer.String())
		rightTopView.SetOrigin(0, 0)
	} else {
		fmt.Fprint(rightTopView, "EMPTY API")
	}
}

// RefreshRequestRecordList 刷新请求记录列表
func RefreshRequestRecordList(g *gocui.Gui) {
	view, _ := g.View("request-history")
	view.Clear()
	list := service.RequestRecordList()
	for _, record := range list {
		if record.Id == entity.SelectedQuestRecord {
		    fmt.Fprintf(view, "> \033[34;1m%s\033[0m [\033[31;1m%s\033[0m] \n",
		                record.RequestTime.Local().Format("2006-01-02 15:04:05"), record.Path)
		} else {
		    fmt.Fprintf(view, "  \033[34;1m%s\033[0m [\033[31;1m%s\033[0m] \n",
		                record.RequestTime.Local().Format("2006-01-02 15:04:05"), record.Path)
		}
	}

    if common.ViewActiveIndex != common.ViewIndexMap["request-history"] {
    	return
    }

	index := slices.IndexFunc(list, func(x entity.RequestRecord) bool {
	    return x.Id == entity.SelectedQuestRecord
	})

	rightTopView, _ := g.View("right-top")
	rightTopView.Clear()
	rightBottomView, _ := g.View("right-bottom")
	rightBottomView.Clear()

	if entity.SelectedQuestRecord != -1 {
		view.SetCursor(0, index)
		api, _ := db.Find(entity.SelectedQuestRecord)
		var buffer bytes.Buffer
	    format_params, _ := utils.PrettyPrintJSON(api.Params)
	    format_respond, _ := utils.PrettyPrintJSON(api.Respond)
		fmt.Fprintf(&buffer, "\033[34;1mName\033[0m: %s \t \033[34;1mMethod\033[0m: %s\n",api.Name, api.Method)
		fmt.Fprintf(&buffer, "\033[34;1mPath\033[0m: %s\n", api.Path)
		fmt.Fprintf(&buffer, "\033[34;1mParams\033[0m: \n%s\n", format_params)
		fmt.Fprint(rightTopView, buffer.String())
		fmt.Fprint(rightBottomView, format_respond)
		rightTopView.SetOrigin(0, 0)
	}
}

func MoveAPISelectionUp(g *gocui.Gui, v *gocui.View) error {
    return MoveSelection(
        service.APIList,
        func(api entity.API) int { return api.Id },
        &entity.SelectedAPI,
        -1,
        func() { UpdateAPIList(g) },
    )
}

func MoveAPISelectionDown(g *gocui.Gui, v *gocui.View) error {
    return MoveSelection(
        service.APIList,
        func(api entity.API) int { return api.Id },
        &entity.SelectedAPI,
        1,
        func() { UpdateAPIList(g) },
    )
}

func MoveRequestRecordSelectionUp(g *gocui.Gui, v *gocui.View) error {
    return MoveSelection(
        service.RequestRecordList,
        func(record entity.RequestRecord) int { return record.Id },
        &entity.SelectedQuestRecord,
        -1,
        func() { RefreshRequestRecordList(g) },
    )
}

func MoveRequestRecordSelectionDown(g *gocui.Gui, v *gocui.View) error {
    return MoveSelection(
        service.RequestRecordList,
        func(record entity.RequestRecord) int { return record.Id },
        &entity.SelectedQuestRecord,
        1,
        func() { RefreshRequestRecordList(g) },
    )
}

func MoveSelection[T any, ID comparable](
    listProvider func() []T,
    idGetter func(T) ID,
    currentSelected *ID,
    direction int,
    updateUI func(),
) error {
    list := listProvider()
    if len(list) <= 1 {
        return nil
    }

    index := slices.IndexFunc(list, func(x T) bool {
        return idGetter(x) == *currentSelected
    })

    // 如果没找到当前选中项，或者在列表边界，则不移动
    if index == -1 {
        if len(list) > 0 {
            // 如果没找到当前选中项但列表不为空，选择第一项
            *currentSelected = idGetter(list[0])
            updateUI()
        }
        return nil
    }

    newIndex := index + direction
    if newIndex < 0 || newIndex >= len(list) {
        return nil
    }

    *currentSelected = idGetter(list[newIndex])
    updateUI()
    return nil
}
