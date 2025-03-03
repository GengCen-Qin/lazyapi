package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	// 视图分布 - 只有三个视图
	viewArr = []string{"left", "right-top", "right-bottom"}
	active  = 0

	// 每个视图对应的状态栏文案
	statusMessages = map[string]string{
		"left":         "接口列表 | n(new), e(edit), d(delete)",
		"right-top":    "接口定义 | s(save), c(cancel)",
		"right-bottom": "响应定义 | r(request), f(format)",
	}

	// 表单字段
	formFields = []string{"name", "path", "method"}
	formLabels = map[string]string{
		"name":   "API 名称",
		"path":   "请求路径",
		"method": "请求方式",
	}
	currentField = 0
	formActive   = false
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	// 当视图变化时更新状态栏
	if err := updateStatusBar(g, name); err != nil {
		return nil, err
	}

 	// 如果是表单相关视图，确保整个表单及其字段都在顶部
    if formActive && strings.HasPrefix(name, "form") {
        // 首先将表单容器置于顶部
        if _, err := g.SetViewOnTop("form"); err != nil {
            return nil, err
        }

        // 将表单按钮置于顶部
        if _, err := g.SetViewOnTop("form-buttons"); err != nil {
            return nil, err
        }

        // 将所有表单字段置于顶部
        for _, field := range formFields {
            fieldName := "form-" + field
            if _, err := g.SetViewOnTop(fieldName); err != nil {
                return nil, err
            }
        }
    }

	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	// 所有视图都是可编辑的
	g.Cursor = true

	active = nextIndex
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// 底部状态栏: 横跨整个宽度
  	if v, err := g.SetView("status", 0, maxY-2, maxX-1, maxY); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Wrap = true
        v.Frame = false  // 不显示边框
        v.FgColor = gocui.ColorBlue // 设置前景色
    }

    // 左边视图占宽度的 3/10
    leftWidth := int(float64(maxX) * 0.3)

    // 左边视图: 占整个高度，宽度为 3/10
    if v, err := g.SetView("left", 0, 0, leftWidth-1, maxY-2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "接口列表"
        v.Wrap = true

        if _, err = setCurrentViewOnTop(g, "left"); err != nil {
            return err
        }
    }

    // 右上视图: 宽度 7/10，高度为 1/2
    if v, err := g.SetView("right-top", leftWidth, 0, maxX-1, maxY/2-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "接口定义"
        v.Wrap = true
        v.Autoscroll = true
    }

    // 右下视图: 宽度 7/10，高度为 1/2
    if v, err := g.SetView("right-bottom", leftWidth, maxY/2, maxX-1, maxY-2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "响应定义"
        v.Editable = true
        v.Wrap = true
        fmt.Fprint(v, "Press TAB to change current view")
    }

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	// 左侧视图键绑定 - 'n'键创建新API
	if err := g.SetKeybinding("left", 'n', gocui.ModNone, showNewAPIForm); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

// 更新状态栏信息
func updateStatusBar(g *gocui.Gui, viewName string) error {
	statusView, err := g.View("status")
	if err != nil {
		return err
	}

	statusView.Clear()
	message, exists := statusMessages[viewName]
	if !exists {
		message = "按 TAB 切换视图 | Ctrl+C 退出"
	}

	fmt.Fprintf(statusView, "%s", message)
	return nil
}

// 修改showNewAPIForm函数，确保字段获得焦点
func showNewAPIForm(g *gocui.Gui, v *gocui.View) error {
    maxX, maxY := g.Size()
    currentField = 0 // 重置当前字段为第一个

    // 创建表单容器
    if v, err := g.SetView("form", maxX/6, maxY/6, maxX*5/6, maxY*5/6); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "新建API"
        v.Wrap = true
        formActive = true
    }

    // 创建表单字段
    for i, field := range formFields {
        label := formLabels[field]
        fieldName := "form-" + field
        fieldView, err := g.SetView(fieldName, maxX/6+1, maxY/6+2+i*3, maxX*5/6-1, maxY/6+4+i*3)
        if err != nil && err != gocui.ErrUnknownView {
            return err
        }
        fieldView.Title = label
        fieldView.Editable = true
        fieldView.Wrap = true
        if field == "method" {
            fmt.Fprint(fieldView, "GET")
        }

        // 为每个字段添加键绑定
        if err := g.SetKeybinding(fieldName, gocui.KeyTab, gocui.ModNone, nextFormField); err != nil {
            return err
        }
        if err := g.SetKeybinding(fieldName, gocui.KeyEnter, gocui.ModNone, saveNewAPI); err != nil {
            return err
        }
        if err := g.SetKeybinding(fieldName, gocui.KeyCtrlQ, gocui.ModNone, closeForm); err != nil {
            return err
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
    for _, field := range formFields {
        g.SetViewOnTop("form-" + field)
    }
    g.SetViewOnTop("form-buttons")

    // 设置初始焦点到第一个字段
    if _, err := setCurrentViewOnTop(g, "form-"+formFields[0]); err != nil {
        return err
    }

    return nil
}

// 关闭表单
func closeForm(g *gocui.Gui, v *gocui.View) error {
    if !formActive {
        return nil
    }

    // 删除所有表单视图
    g.DeleteView("form")
    for _, field := range formFields {
        fieldName := "form-" + field
        g.DeleteView(fieldName)

        // 删除各个字段的键绑定
        g.DeleteKeybinding(fieldName, gocui.KeyEnter, gocui.ModNone)
        g.DeleteKeybinding(fieldName, gocui.KeyEsc, gocui.ModNone)
        g.DeleteKeybinding(fieldName, gocui.KeyTab, gocui.ModNone)
    }
    g.DeleteView("form-buttons")

    // 重新设置焦点到左侧视图
    if _, err := setCurrentViewOnTop(g, "left"); err != nil {
        return err
    }
    formActive = false

    return nil
}

// 保存新API
func saveNewAPI(g *gocui.Gui, v *gocui.View) error {
	if !formActive {
		return nil
	}

	// 收集表单数据
	var name, path, method string
	nameView, _ := g.View("form-name")
	pathView, _ := g.View("form-path")
	methodView, _ := g.View("form-method")

	name = strings.TrimSpace(nameView.Buffer())
	path = strings.TrimSpace(pathView.Buffer())
	method = strings.TrimSpace(methodView.Buffer())

	// 简单验证
	if name == "" || path == "" || method == "" {
		statusView, _ := g.View("status")
		statusView.Clear()
		fmt.Fprint(statusView, "错误: 所有字段必须填写")
		return nil
	}

	// 添加新API到左侧列表
	leftView, _ := g.View("left")
	fmt.Fprintf(leftView, "%s [%s] %s\n", name, method, path)

	// 展示API定义在右上视图
	rightTopView, _ := g.View("right-top")
	rightTopView.Clear()
	fmt.Fprintf(rightTopView, "API名称: %s\n", name)
	fmt.Fprintf(rightTopView, "请求路径: %s\n", path)
	fmt.Fprintf(rightTopView, "请求方式: %s\n", method)

	// 关闭表单
	return closeForm(g, v)
}

// 在表单字段间切换
func nextFormField(g *gocui.Gui, v *gocui.View) error {
	nextField := (currentField + 1) % len(formFields)
	fieldName := "form-" + formFields[nextField]

	if _, err := setCurrentViewOnTop(g, fieldName); err != nil {
		return err
	}

	currentField = nextField
	return nil
}
