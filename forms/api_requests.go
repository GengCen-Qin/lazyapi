package forms

import (
	"encoding/json"
	"fmt"
	"strings"

	"lazyapi/common"
	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/utils"

	"github.com/GengCen-Qin/gocui"
	"github.com/go-resty/resty/v2"
)

// RequestAPI 请求当前选中的API
func RequestAPI(g *gocui.Gui, v *gocui.View) error {
	if entity.SelectedAPI == -1 {
		return nil
	}

	api, _ := db.FindAPI(entity.SelectedAPI)
	params, err := api.GetParams()
	if err != nil {
		return err
	}

	if len(params) == 0 {
		return sendRequest(g, api, params)
	}

	// 创建新的view来展示和编辑params
	maxX, maxY := g.Size()
	viewWidth, viewHeight := 50, 8
	left := maxX/2 - viewWidth/2
	right := maxX/2 + viewWidth/2
	top := maxY/2 - viewHeight/2
	bottom := maxY/2 + viewHeight/2

	pv, err := g.SetView("requestConfirmView", left, top, right, bottom)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		pv.Title = "确认请求参数(可编辑)"
		pv.Editable = true
		pv.Wrap = true
		g.SetCurrentView("requestConfirmView")
		g.Cursor = true

		status_view, _ := g.View("status")
		status_view.Clear()
		fmt.Fprint(status_view, common.StatusMessages["requestConfirmView"])

		// 在 view 中显示默认的 params
		fmt.Fprintln(pv, api.Params)

		g.SetKeybinding("requestConfirmView", gocui.KeyCtrlR, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			// 从 view 获取输入的 params
			pv, _ := g.View("requestConfirmView")
			inputParams := pv.Buffer()
			var inputParamsMap map[string]interface{}
			jsonErr := json.Unmarshal([]byte(inputParams), &inputParamsMap)
			if jsonErr != nil {
				statusView, _ := g.View("status")
				statusView.Clear()
				fmt.Fprintf(statusView, "\033[31;1m%s\033[0m", "params must be a json!!!")
				return nil
			}
			// 发送请求
			err := sendRequest(g, api, inputParamsMap)

			// Remove the requestConfirmView and associated keybindings after the request
			if err := g.DeleteView("requestConfirmView"); err != nil {
				return err
			}

			// Remove keybindings
			if err := g.DeleteKeybinding("requestConfirmView", gocui.KeyCtrlR, gocui.ModNone); err != nil {
				return err
        	}

            g.Cursor = false

            status_view.Clear()
            fmt.Fprint(status_view, "request success")

			return err
		})

		g.SetKeybinding("requestConfirmView", gocui.KeyCtrlQ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			status_view.Clear()
            fmt.Fprint(status_view, "request cancel")
			g.DeleteView("requestConfirmView")
		 	g.DeleteKeybinding("requestConfirmView", gocui.KeyCtrlQ, gocui.ModNone)
            g.Cursor = false
			return nil
		})
	}

	return nil
}

// mapToStringMap 将interface{}类型的map转换为string类型的map
func mapToStringMap(params map[string]interface{}) map[string]string {
    stringMap := make(map[string]string)

    for k, v := range params {
        stringMap[k] = fmt.Sprint(v)
    }

    return stringMap
}

// MapToJSONString 将map转换为JSON字符串
func MapToJSONString(params map[string]interface{}) (string, error) {
    // Marshal the map into JSON bytes
    jsonBytes, err := json.Marshal(params)
    if err != nil {
        return "", fmt.Errorf("error marshaling map to JSON: %w", err)
    }

    // Convert bytes to string
    jsonString := string(jsonBytes)

    return jsonString, nil
}

// sendRequest 发送API请求
func sendRequest(g *gocui.Gui, api *entity.API, params map[string]interface{}) error {
	client := resty.New()

	request := client.R().
    SetHeader("Content-Type", "application/json")

	method := strings.ToUpper(api.Method)
	if method == "GET" {
	    // 对于GET请求，将参数设置为查询参数
	    request.SetQueryParams(mapToStringMap(params))
	} else {
	    // 对于其他请求(POST, PUT等)，将参数设置为请求体
	    request.SetBody(params)
	}

	resp, err := request.Execute(method, api.Path)

	bottomView, _ := g.View("right-bottom")
	bottomView.Clear()
	utils.ResetViewOrigin(bottomView)

	json_params, _ := MapToJSONString(params)

	if err != nil {
		db.InsertRequestRecord(api, json_params, err.Error())
		fmt.Fprint(bottomView, "请求失败: ", err)
	} else {
		respBody := resp.Body()
		format_josn, _ := utils.PrettyPrintJSON(string(respBody))
		db.InsertRequestRecord(api, json_params, string(respBody))
		fmt.Fprint(bottomView, format_josn)
	}

	RefreshRequestRecordList(g)
	g.SetCurrentView("left")
	return nil
}
