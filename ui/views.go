package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"lazyapi/common"
)

// SetCurrentViewOnTop 设置当前视图为顶层并更新状态栏
func SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	// 当视图变化时更新状态栏
	if err := UpdateStatusBar(g, name); err != nil {
		return nil, err
	}

	// 如果是表单相关视图，确保整个表单及其字段都在顶部
	if common.FormInfo.Active && strings.HasPrefix(name, "form") {
		// 首先将表单容器置于顶部
		if _, err := g.SetViewOnTop("form"); err != nil {
			return nil, err
		}

		// 将表单按钮置于顶部
		if _, err := g.SetViewOnTop("form-buttons"); err != nil {
			return nil, err
		}

		// 将所有表单字段置于顶部
		for _, field := range common.FormInfo.Fields {
			fieldName := "form-" + field
			if _, err := g.SetViewOnTop(fieldName); err != nil {
				return nil, err
			}
		}
	}

	return g.SetViewOnTop(name)
}

// NextView 切换到下一个视图
func NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (common.Active + 1) % len(common.ViewArr)
	name := common.ViewArr[nextIndex]

	if _, err := SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	common.Active = nextIndex
	return nil
}

// UpdateStatusBar 更新状态栏信息
func UpdateStatusBar(g *gocui.Gui, viewName string) error {
	statusView, err := g.View("status")
	if err != nil {
		return err
	}

	statusView.Clear()
	message, exists := common.StatusMessages[viewName]
	if !exists {
		message = "按 TAB 切换视图 | Ctrl+C 退出"
	}

	fmt.Fprintf(statusView, "%s", message)
	return nil
}
