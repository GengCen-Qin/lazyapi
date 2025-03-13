package forms

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"lazyapi/common"
	"lazyapi/models"
	"lazyapi/ui"

	"github.com/go-resty/resty/v2"
	"github.com/GengCen-Qin/gocui"
)

// ShowNewAPIForm 显示新建API表单
func ShowNewAPIForm(g *gocui.Gui, v *gocui.View) error {
	common.FormInfo.Active = true
	maxX, maxY := g.Size()
	common.FormInfo.CurrentField = 0 // 重置当前字段为第一个

	// 创建表单容器
	if v, err := g.SetView("form", maxX/6, maxY/6, maxX*5/6, maxY*5/6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
	}

	var form_view, _ = g.View("form")
	if common.FormInfo.IsEditing {
    	form_view.Title = "编辑API"
	} else {
		form_view.Title = "新建API"
	}

	// 创建表单字段
	for _, field := range common.FormInfo.Fields {
		label := common.FormInfo.Labels[field]
		fieldName := "form-" + field
		var fieldView *gocui.View
		var err error

		// 调整视图布局
		switch field {
		case "name", "method":
			// 将 "名称" 和 "请求方式" 放在同一行
			if field == "name" {
				fieldView, err = g.SetView(fieldName, maxX/6+1, maxY/6+2, maxX/2-1, maxY/6+4)
			} else {
				fieldView, err = g.SetView(fieldName, maxX/2+1, maxY/6+2, maxX*5/6-1, maxY/6+4)
			}
		case "path":
			// "请求路径" 放在第二行
			fieldView, err = g.SetView(fieldName, maxX/6+1, maxY/6+5, maxX*5/6-1, maxY/6+7)
		case "params":
			// "请求参数" 放在第三行，并且视图更大
			fieldView, err = g.SetView(fieldName, maxX/6+1, maxY/6+8, maxX*5/6-1, maxY*5/6-1)
		}

		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		fieldView.Title = label
		fieldView.Editable = true
		fieldView.Wrap = true
		if field == "method" {
			fmt.Fprint(fieldView, "GET")
		}
		if field == "params" {
			fmt.Fprint(fieldView, "{}")
		}
		if field == "path" {
			fmt.Fprint(fieldView, "http://")
		}

		// 为每个字段添加键绑定
		if err := g.SetKeybinding(fieldName, gocui.KeyTab, gocui.ModNone, NextFormField); err != nil {
			return err
		}
		if err := g.SetKeybinding(fieldName, gocui.KeyEsc, gocui.ModNone, CloseForm); err != nil {
			return err
		}
		if field != "params" {
			if err := g.SetKeybinding(fieldName, gocui.KeyEnter, gocui.ModNone, SaveNewAPI); err != nil {
				return err
			}
			if err := g.SetKeybinding(fieldName, gocui.KeyArrowDown, gocui.ModNone, NextFormField); err != nil {
				return err
			}
			if err := g.SetKeybinding(fieldName, gocui.KeyArrowUp, gocui.ModNone, BeforeFormField); err != nil {
				return err
			}
		}
	}

	// 添加按钮
	if v, err := g.SetView("form-buttons", maxX/6+1, maxY*5/6-3, maxX*5/6-1, maxY*5/6-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprint(v, "保存(Enter)  取消(Esc)")
	}

	// 确保表单及其所有字段保持在最顶层
	g.SetViewOnTop("form")
	for _, field := range common.FormInfo.Fields {
		g.SetViewOnTop("form-" + field)
	}
	g.SetViewOnTop("form-buttons")

	// 设置初始焦点到第一个字段
	if _, err := ui.SetCurrentViewOnTop(g, "form-"+common.FormInfo.Fields[0]); err != nil {
		return err
	}

	g.Cursor = true

	return nil
}

// CloseForm 关闭表单
func CloseForm(g *gocui.Gui, v *gocui.View) error {
	if !common.FormInfo.Active {
		return nil
	}

	// 删除所有表单视图
	g.DeleteView("form")
	for _, field := range common.FormInfo.Fields {
		fieldName := "form-" + field
		g.DeleteView(fieldName)

		// 删除各个字段的键绑定
		g.DeleteKeybinding(fieldName, gocui.KeyEnter, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyEsc, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyTab, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyArrowUp, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyArrowDown, gocui.ModNone)
	}
	g.DeleteView("form-buttons")

	// 重新设置焦点到左侧视图
	if _, err := ui.SetCurrentViewOnTop(g, "left"); err != nil {
		return err
	}
	common.Active = 0
	common.FormInfo.Active = false
	g.Cursor = false
	return nil
}

// SaveNewAPI 保存API（新建或编辑）
func SaveNewAPI(g *gocui.Gui, v *gocui.View) error {
	if !common.FormInfo.Active {
		return nil
	}
	if common.FormInfo.IsDelete {
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
	if common.FormInfo.IsEditing {
		api, _ := models.FindAPI(models.SelectedAPI)
		api.Name = name
		api.Path = path
		api.Method = method
		api.Params = params
		models.UpdateAPI(api)
	} else {
		// 否则，创建新API并添加到列表
		newAPI := models.NewAPI(name, path, method, params)
		models.SelectedAPI = newAPI.Id
	}

	// 更新视图
	UpdateAPIList(g)

	// 关闭表单
	common.FormInfo.IsEditing = false
	// 关闭光标
	g.Cursor = false
	return CloseForm(g, v)
}

func validateAPIForm(g *gocui.Gui, name, path, method, params string) error {
	var emptyFields []string
	if name == "" {
	    emptyFields = append(emptyFields, "name")
	}
	if path == "" {
	    emptyFields = append(emptyFields, "path")
	}
	if method == "" {
	    emptyFields = append(emptyFields, "method")
	}

	if len(emptyFields) > 0 {
	    warnFormItem(g, emptyFields)
	    return fmt.Errorf("the following fields must be filled: %v", emptyFields)
	}

	if !(strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")) {
  		warnFormItem(g, []string{"path"})
		emptyFields = append(emptyFields, "path")
	}

    // 验证 params 是否为有效的 JSON
    if !json.Valid([]byte(params)) {
   		warnFormItem(g, []string{"params"})
        return fmt.Errorf("params must be a json")
    }

    // 验证 method 是否为 GET 或 POST
    method = strings.ToUpper(method)
    if method != "GET" && method != "POST" {
  		warnFormItem(g, []string{"method"})
        return fmt.Errorf("method must be GET or POST !!!")
    }

    return nil
}

func warnFormItem(g *gocui.Gui, fields []string) {
	for _, view := range fields {
		view, _ := g.View("form-"+view)
		view.BgColor = gocui.ColorRed
    }
}

// NextFormField 在表单字段间切换
func NextFormField(g *gocui.Gui, v *gocui.View) error {
	nextField := (common.FormInfo.CurrentField + 1) % len(common.FormInfo.Fields)
	fieldName := "form-" + common.FormInfo.Fields[nextField]

	if _, err := ui.SetCurrentViewOnTop(g, fieldName); err != nil {
		return err
	}

	common.FormInfo.CurrentField = nextField
	return nil
}

func BeforeFormField(g *gocui.Gui, v *gocui.View) error {
	// 计算前一个字段的索引
	prevField := (common.FormInfo.CurrentField - 1 + len(common.FormInfo.Fields)) % len(common.FormInfo.Fields)
	fieldName := "form-" + common.FormInfo.Fields[prevField]

	// 将焦点设置到前一个字段
	if _, err := ui.SetCurrentViewOnTop(g, fieldName); err != nil {
		return err
	}

	// 更新当前字段索引
	common.FormInfo.CurrentField = prevField
	return nil
}

// SetupFormKeybindings 为表单设置键绑定
func SetupFormKeybindings(g *gocui.Gui) error {
	// 左侧视图键绑定 - 'n'键创建新API或取消删除
	if err := g.SetKeybinding("left", 'n', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if common.FormInfo.IsDelete {
			return CancelDeleteAPI(g, v)
		}
		return ShowNewAPIForm(g, v)
	}); err != nil {
		return err
	}

	// 左侧视图键绑定 - 'e'键编辑选中的API
	if err := g.SetKeybinding("left", 'e', gocui.ModNone, EditAPIForm); err != nil {
		return err
	}

	// 左侧视图键绑定 - 'd'键删除选中的API
	if err := g.SetKeybinding("left", 'd', gocui.ModNone, DeleteAPI); err != nil {
		return err
	}

	// 列表选中API，跳入详情View
	if err := g.SetKeybinding("left", gocui.KeySpace, gocui.ModNone, JumpDetailView); err != nil {
		return err
	}

	// 右上视图键绑定 - 'r'键请求当前API
	if err := g.SetKeybinding("left", 'r', gocui.ModNone, RequestAPI); err != nil {
		return err
	}

	// 添加上下键绑定
	if err := g.SetKeybinding("left", gocui.KeyArrowUp, gocui.ModNone, MoveSelectionUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("left", gocui.KeyArrowDown, gocui.ModNone, MoveSelectionDown); err != nil {
		return err
	}

	return nil
}

// MoveSelectionUp 向上移动选择
func MoveSelectionUp(g *gocui.Gui, v *gocui.View) error {
	list := models.APIList()
	if len(list) == 0 || len(list) == 1 || models.SelectedAPI <= 0 {
		return nil
	}

	index := slices.IndexFunc(list, func(x models.API) bool {
	    return x.Id == models.SelectedAPI
	})

	if index != 0 {
		models.SelectedAPI = list[index-1].Id
	}
	UpdateAPIList(g)
	return nil
}

// MoveSelectionDown 向下移动选择
func MoveSelectionDown(g *gocui.Gui, v *gocui.View) error {
	list := models.APIList()
	if len(list) == 0 || len(list) == 1 {
		return nil
	}

	index := slices.IndexFunc(list, func(x models.API) bool {
	    return x.Id == models.SelectedAPI
	})

	if index + 1 >= len(list) {
		return nil
	}

	models.SelectedAPI = list[index+1].Id

	UpdateAPIList(g)
	return nil
}

// UpdateAPIList 更新左侧API列表显示
func UpdateAPIList(g *gocui.Gui) {
	leftView, leftViewError := g.View("left")
	if leftViewError != nil {
		return // 如果视图不存在，直接返回
	}
	rightTopView, _ := g.View("right-top")
	list := models.APIList()
	leftView.Clear()
	for _, api := range list {
		if api.Id == models.SelectedAPI {
			fmt.Fprintf(leftView, "> %s [%s] \n", api.Name, api.Method)
		} else {
			// 文字颜色控制
			fmt.Fprintf(leftView, "  %s [\033[31;1m%s\033[0m] \n", api.Name, api.Method)
		}
	}

	index := slices.IndexFunc(list, func(x models.API) bool {
	    return x.Id == models.SelectedAPI
	})

	rightTopView.Clear()

	if models.SelectedAPI != -1 {
		leftView.SetCursor(0, index)
		api, _ := models.FindAPI(models.SelectedAPI)
		fmt.Fprintf(rightTopView, "\033[34;1mName\033[0m: %s \t \033[34;1mMethod\033[0m: %s\n",api.Name, api.Method)
		fmt.Fprintf(rightTopView, "\033[34;1mPath\033[0m: %s\n", api.Path)
		fmt.Fprintf(rightTopView, "\033[34;1mParams\033[0m: \n%s\n", api.Params)
	} else {
		fmt.Fprint(rightTopView, "EMPTY API")
	}
}

func RefreshRequestRecordList(g *gocui.Gui) {
	view, _ := g.View("request-history")

	list := models.RequestRecordList()

	view.Clear()
	for _, record := range list {
		if record.Id == models.SelectedQuestRecord {
		    fmt.Fprintf(view, ">\033[34;1m%s\033[0m [\033[a31;1m%s\033[0m] \n",
		                record.RequestTime.Local().Format("2006-01-02 15:04:05"), record.Path)
		} else {
		    fmt.Fprintf(view, " \033[34;1m%s\033[0m [\033[31;1m%s\033[0m] \n",
		                record.RequestTime.Local().Format("2006-01-02 15:04:05"), record.Path)
		}
	}
}

func EditAPIForm(g *gocui.Gui, v *gocui.View) error {
	if models.SelectedAPI == -1 {
		return nil
	}

	// 标记为编辑模式
	common.FormInfo.IsEditing = true

	// 显示表单并填充数据
	if err := ShowNewAPIForm(g, v); err != nil {
		return err
	}

	// 获取当前选中的API
	api, _ := models.FindAPI(models.SelectedAPI)

	fillFormFields(g, api)
	return nil
}

func fillFormFields(g *gocui.Gui, api * models.API) error {
    // 定义要填充的字段映射
    fields := map[string]string{
        "form-name":   api.Name,
        "form-path":   api.Path,
        "form-method": api.Method,
        "form-params": api.Params,
    }

    // 统一处理所有字段
    for viewName, value := range fields {
        view, err := g.View(viewName)
        if err != nil {
            continue // 或者返回错误：return fmt.Errorf("无法获取视图 %s: %v", viewName, err)
        }
        view.Clear()
        fmt.Fprint(view, value)
    }

    return nil
}

// DeleteAPI 删除选中的API
func DeleteAPI(g *gocui.Gui, v *gocui.View) error {
	if models.SelectedAPI == -1 {
		return nil
	}

	common.FormInfo.IsDelete = true

	// 显示确认提示
	statusView, _ := g.View("status")
	statusView.Clear()
	fmt.Fprint(statusView, "confirm to delete ? (y/n)")

	// 绑定确认和取消操作
	if err := g.SetKeybinding("", 'y', gocui.ModNone, ConfirmDeleteAPI); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'n', gocui.ModNone, CancelDeleteAPI); err != nil {
		return err
	}

	return nil
}

func JumpDetailView(g *gocui.Gui, v *gocui.View) error {
    index := slices.IndexFunc(common.ViewArr, func(x string) bool {
		return x == "right-top"
	})
	if _, err := ui.SetCurrentViewOnTop(g, "right-top"); err != nil {
		return err
	}
	common.Active = index
	return nil
}

func RequestAPI(g *gocui.Gui, v *gocui.View) error {
	if models.SelectedAPI == -1 {
		return nil
	}

	api, _ := models.FindAPI(models.SelectedAPI)
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
				return jsonErr // handle error appropriately
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

func mapToStringMap(params map[string]interface{}) map[string]string {
    stringMap := make(map[string]string)

    for k, v := range params {
        stringMap[k] = fmt.Sprint(v)
    }

    return stringMap
}

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

func sendRequest(g *gocui.Gui, api *models.API, params map[string]interface{}) error {
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

	json_params, _ := MapToJSONString(params)

	if err != nil {
		models.InsertRequestRecord(api, json_params, err.Error())
		fmt.Fprint(bottomView, "请求失败: ", err)
	} else {
		respBody := resp.Body()
		models.InsertRequestRecord(api, json_params, string(respBody))
		fmt.Fprint(bottomView, string(respBody))
	}

	RefreshRequestRecordList(g)
	return nil
}

// ConfirmDeleteAPI 确认删除API
func ConfirmDeleteAPI(g *gocui.Gui, v *gocui.View) error {
	models.DeleteAPI(models.SelectedAPI)

    list := models.APIList()
	if len(list) == 0 {
		models.SelectedAPI = -1
	} else {
		models.SelectedAPI = list[len(list)-1].Id
	}

	// 更新视图
	UpdateAPIList(g)

	// 清除确认提示
	statusView, _ := g.View("status")
	statusView.Clear()
	fmt.Fprint(statusView, "delete success !!!")

	// 删除临时键绑定
	g.DeleteKeybinding("", 'y', gocui.ModNone)
	g.DeleteKeybinding("", 'n', gocui.ModNone)

	// 重置删除标志
	common.FormInfo.IsDelete = false

	return nil
}

// CancelDeleteAPI 取消删除API
func CancelDeleteAPI(g *gocui.Gui, v *gocui.View) error {
	// 清除确认提示
	statusView, _ := g.View("status")
	statusView.Clear()

	// 删除临时键绑定
	g.DeleteKeybinding("", 'y', gocui.ModNone)
	g.DeleteKeybinding("", 'n', gocui.ModNone)

	// 将焦点重新设置到 left 视图
	if _, err := ui.SetCurrentViewOnTop(g, "left"); err != nil {
		return err
	}

	// 重置删除标志
	common.FormInfo.IsDelete = false

	return nil
}
