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
	maxX, maxY := g.Size()
	FormInfo.CurrentField = 0 // 重置当前字段为第一个

	// 创建表单容器
	if v, err := g.SetView("form", maxX/6, maxY/6, maxX*5/6, maxY*5/6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
	}

	var form_view, _ = g.View("form")
	if FormInfo.IsEditing {
    	form_view.Title = "编辑API"
	} else {
		form_view.Title = "新建API"
	}

	// 创建表单字段
	for _, field := range FormInfo.Fields {
		label := FormInfo.Labels[field]
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
	for _, field := range FormInfo.Fields {
		g.SetViewOnTop("form-" + field)
	}
	g.SetViewOnTop("form-buttons")

	// 设置初始焦点到第一个字段
	if _, err := SetCurrentViewOnTop(g, "form-"+FormInfo.Fields[0]); err != nil {
		return err
	}

	g.Cursor = true

	return nil
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
	g.DeleteView("form-buttons")

	// 重新设置焦点到左侧视图
	if _, err := SetCurrentViewOnTop(g, "left"); err != nil {
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
