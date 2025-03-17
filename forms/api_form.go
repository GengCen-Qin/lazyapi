package forms

import (
	"fmt"
	"strings"
	"encoding/json"

	"lazyapi/common"
	"lazyapi/models/entity"

	"github.com/GengCen-Qin/gocui"
)

// ShowNewAPIForm 显示新建API表单
func ShowNewAPIForm(g *gocui.Gui, v *gocui.View) error {
    FormInfo.Active = true
    FormInfo.CurrentField = 0

    if err := createFormContainer(g); err != nil {
        return err
    }

    if err := setupFormTitle(g); err != nil {
        return err
    }

    if err := createFormFields(g); err != nil {
        return err
    }

    if err := setupViewOrder(g); err != nil {
        return err
    }

    g.Cursor = true
    return nil
}

// createFormContainer 创建表单容器
func createFormContainer(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("form", maxX/6, maxY/6, maxX*5/6, maxY*5/6); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Wrap = true
    }
    return nil
}

// setupFormTitle 设置表单标题
func setupFormTitle(g *gocui.Gui) error {
    form_view, err := g.View("form")
    if err != nil {
        return err
    }

    if FormInfo.IsEditing {
        form_view.Title = "编辑API"
    } else {
        form_view.Title = "新建API"
    }
    return nil
}

// getFieldPosition 获取字段位置
func getFieldPosition(field string, maxX, maxY int) (x1, y1, x2, y2 int) {
    switch field {
    case "name":
        return maxX/6+1, maxY/6+2, maxX/2-1, maxY/6+4
    case "method":
        return maxX/2+1, maxY/6+2, maxX*5/6-1, maxY/6+4
    case "path":
        return maxX/6+1, maxY/6+5, maxX*5/6-1, maxY/6+7
    case "params":
        return maxX/6+1, maxY/6+8, maxX*5/6-1, maxY*5/6-1
    default:
        return 0, 0, 0, 0
    }
}

// createFormField 创建单个表单字段
func createFormField(g *gocui.Gui, field string) error {
    maxX, maxY := g.Size()
    fieldName := "form-" + field

    x1, y1, x2, y2 := getFieldPosition(field, maxX, maxY)
    fieldView, err := g.SetView(fieldName, x1, y1, x2, y2)
    if err != nil && err != gocui.ErrUnknownView {
        return err
    }

    setupFieldProperties(fieldView, field)
    return setupFieldKeybindings(g, fieldName, field)
}

// setupFieldProperties 设置字段属性
func setupFieldProperties(v *gocui.View, field string) {
    v.Title = FormInfo.Labels[field]
    v.Editable = true
    v.Wrap = true

    // 设置默认值
    defaultValues := map[string]string{
        "method": "GET",
        "params": "{}",
        "path":   "http://",
    }

    if defaultValue, exists := defaultValues[field]; exists {
        fmt.Fprint(v, defaultValue)
    }
}

// createFormFields 创建所有表单字段
func createFormFields(g *gocui.Gui) error {
    for _, field := range FormInfo.Fields {
        if err := createFormField(g, field); err != nil {
            return err
        }
    }
    return nil
}

// setupFieldKeybindings 设置字段键绑定
func setupFieldKeybindings(g *gocui.Gui, fieldName, field string) error {
    keybindings := []struct {
        key interface{}
        handler func(*gocui.Gui, *gocui.View) error
    }{
        {gocui.KeyTab, NextFormField},
        {gocui.KeyEsc, CloseForm},
    }

    if field != "params" {
        keybindings = append(keybindings,
            struct {
                key interface{}
                handler func(*gocui.Gui, *gocui.View) error
            }{
                gocui.KeyEnter, SaveNewAPI,
            },
            struct {
                key interface{}
                handler func(*gocui.Gui, *gocui.View) error
            }{
                gocui.KeyArrowDown, NextFormField,
            },
            struct {
                key interface{}
                handler func(*gocui.Gui, *gocui.View) error
            }{
                gocui.KeyArrowUp, BeforeFormField,
            },
        )
    }

    for _, kb := range keybindings {
        if err := g.SetKeybinding(fieldName, kb.key, gocui.ModNone, kb.handler); err != nil {
            return err
        }
    }
    return nil
}

// setupViewOrder 设置视图顺序
func setupViewOrder(g *gocui.Gui) error {
    g.SetViewOnTop("form")
    for _, field := range FormInfo.Fields {
        g.SetViewOnTop("form-" + field)
    }

    // 设置初始焦点
    _, err := SetCurrentViewOnTop(g, "form-"+FormInfo.Fields[0])
    return err
}

// CloseForm 关闭表单
func CloseForm(g *gocui.Gui, v *gocui.View) error {
	if !FormInfo.Active {
		return nil
	}

	// 删除所有表单视图
	g.DeleteView("form")
	for _, field := range FormInfo.Fields {
		fieldName := "form-" + field
		g.DeleteView(fieldName)

		// 删除各个字段的键绑定
		g.DeleteKeybinding(fieldName, gocui.KeyEnter, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyEsc, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyTab, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyArrowUp, gocui.ModNone)
		g.DeleteKeybinding(fieldName, gocui.KeyArrowDown, gocui.ModNone)
	}

	// 重新设置焦点到左侧视图
	if _, err := SetCurrentViewOnTop(g, "api_list"); err != nil {
		return err
	}
	common.ViewActiveIndex = 0
	FormInfo.Active = false
	g.Cursor = false
	return nil
}

// NextFormField 在表单字段间切换
func NextFormField(g *gocui.Gui, v *gocui.View) error {
	nextField := (FormInfo.CurrentField + 1) % len(FormInfo.Fields)
	fieldName := "form-" + FormInfo.Fields[nextField]

	if _, err := SetCurrentViewOnTop(g, fieldName); err != nil {
		return err
	}

	FormInfo.CurrentField = nextField
	return nil
}

// BeforeFormField 切换到前一个表单字段
func BeforeFormField(g *gocui.Gui, v *gocui.View) error {
	// 计算前一个字段的索引
	prevField := (FormInfo.CurrentField - 1 + len(FormInfo.Fields)) % len(FormInfo.Fields)
	fieldName := "form-" + FormInfo.Fields[prevField]

	// 将焦点设置到前一个字段
	if _, err := SetCurrentViewOnTop(g, fieldName); err != nil {
		return err
	}

	// 更新当前字段索引
	FormInfo.CurrentField = prevField
	return nil
}

// validateAPIForm 验证表单字段
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

// warnFormItem 高亮显示错误字段
func warnFormItem(g *gocui.Gui, fields []string) {
	for _, view := range fields {
		view, _ := g.View("form-"+view)
		view.BgColor = gocui.ColorRed
    }
}

// fillFormFields 填充表单字段
func fillFormFields(g *gocui.Gui, api *entity.API) error {
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
            continue
        }
        view.Clear()
        fmt.Fprint(view, value)
    }

    return nil
}
