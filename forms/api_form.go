package forms

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"lazyapi/common"
	"lazyapi/models"
	"lazyapi/ui"
)

// ShowNewAPIForm 显示新建API表单
func ShowNewAPIForm(g *gocui.Gui, v *gocui.View) error {

	maxX, maxY := g.Size()
	common.FormInfo.CurrentField = 0 // 重置当前字段为第一个

	// 创建表单容器
	if v, err := g.SetView("form", maxX/6, maxY/6, maxX*5/6, maxY*5/6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if common.FormInfo.IsEditing {
			v.Title = "编辑API"
		} else {
			v.Title = "新建API"
		}

		v.Wrap = true
		common.FormInfo.Active = true
	}

	// 创建表单字段
	for i, field := range common.FormInfo.Fields {
		label := common.FormInfo.Labels[field]
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
		if err := g.SetKeybinding(fieldName, gocui.KeyTab, gocui.ModNone, NextFormField); err != nil {
			return err
		}
		if err := g.SetKeybinding(fieldName, gocui.KeyEnter, gocui.ModNone, SaveNewAPI); err != nil {
			return err
		}
		if err := g.SetKeybinding(fieldName, gocui.KeyCtrlQ, gocui.ModNone, CloseForm); err != nil {
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
	for _, field := range common.FormInfo.Fields {
		g.SetViewOnTop("form-" + field)
	}
	g.SetViewOnTop("form-buttons")

	// 设置初始焦点到第一个字段
	if _, err := ui.SetCurrentViewOnTop(g, "form-"+common.FormInfo.Fields[0]); err != nil {
		return err
	}

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
	}
	g.DeleteView("form-buttons")

	// 重新设置焦点到左侧视图
	if _, err := ui.SetCurrentViewOnTop(g, "left"); err != nil {
		return err
	}
	common.FormInfo.Active = false

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
		fmt.Fprint(statusView, "all fields must input !!!")
		return nil
	}

	// 如果是编辑模式，更新现有API
	if common.FormInfo.IsEditing {
		api := models.APIList[models.SelectedAPI]
		api.Name = name
		api.Path = path
		api.Method = method
	} else {
		// 否则，创建新API并添加到列表
		newAPI := models.NewAPI(name, path, method)
		models.APIList = append(models.APIList, newAPI)
		models.SelectedAPI = len(models.APIList) - 1
	}

	// 更新视图
	UpdateAPIList(g)

	// 关闭表单
	common.FormInfo.IsEditing = false
	return CloseForm(g, v)
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
	if len(models.APIList) == 0 {
		return nil
	}

	if models.SelectedAPI > 0 {
		models.SelectedAPI--
	}
	UpdateAPIList(g)
	return nil
}

// MoveSelectionDown 向下移动选择
func MoveSelectionDown(g *gocui.Gui, v *gocui.View) error {
	if len(models.APIList) == 0 {
		return nil
	}

	if models.SelectedAPI < len(models.APIList)-1 {
		models.SelectedAPI++
	}
	UpdateAPIList(g)
	return nil
}

// UpdateAPIList 更新左侧API列表显示
func UpdateAPIList(g *gocui.Gui) {
	leftView, _ := g.View("left")
	rightTopView, _ := g.View("right-top")

	// 清空并重新渲染左侧视图
	leftView.Clear()
	for i, api := range models.APIList {
		if i == models.SelectedAPI {
			fmt.Fprintf(leftView, "> %s [%s] %s\n", api.Name, api.Method, api.Path)
		} else {
			fmt.Fprintf(leftView, "  %s [%s] %s\n", api.Name, api.Method, api.Path)
		}
	}

	// 设置光标位置
	if models.SelectedAPI >= 0 {
		leftView.SetCursor(0, models.SelectedAPI)
	}

	// 更新右上视图的内容
	rightTopView.Clear()
	if models.SelectedAPI >= 0 && models.SelectedAPI < len(models.APIList) {
		api := models.APIList[models.SelectedAPI]
		fmt.Fprintf(rightTopView, "API名称: %s\n", api.Name)
		fmt.Fprintf(rightTopView, "请求路径: %s\n", api.Path)
		fmt.Fprintf(rightTopView, "请求方式: %s\n", api.Method)
	} else {
		fmt.Fprint(rightTopView, "无选中API")
	}
}

func EditAPIForm(g *gocui.Gui, v *gocui.View) error {
	if len(models.APIList) == 0 || models.SelectedAPI < 0 || models.SelectedAPI >= len(models.APIList) {
		return nil
	}

	// 获取当前选中的API
	api := models.APIList[models.SelectedAPI]

	// 显示表单并填充数据
	if err := ShowNewAPIForm(g, v); err != nil {
		return err
	}

	// 填充表单字段
	if nameView, err := g.View("form-name"); err == nil {
		nameView.Clear()
		fmt.Fprint(nameView, api.Name)
	}
	if pathView, err := g.View("form-path"); err == nil {
		pathView.Clear()
		fmt.Fprint(pathView, api.Path)
	}
	if methodView, err := g.View("form-method"); err == nil {
		methodView.Clear()
		fmt.Fprint(methodView, api.Method)
	}

	// 标记为编辑模式
	common.FormInfo.IsEditing = true
	return nil
}

// DeleteAPI 删除选中的API
func DeleteAPI(g *gocui.Gui, v *gocui.View) error {
	if len(models.APIList) == 0 || models.SelectedAPI < 0 || models.SelectedAPI >= len(models.APIList) {
		return nil
	}

	common.FormInfo.IsEditing = true
	common.FormInfo.IsDelete = true

	// 显示确认提示
	statusView, _ := g.View("status")
	statusView.Clear()
	fmt.Fprint(statusView, "confirm to delete？(y/n)")

	// 绑定确认和取消操作
	if err := g.SetKeybinding("", 'y', gocui.ModNone, ConfirmDeleteAPI); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'n', gocui.ModNone, CancelDeleteAPI); err != nil {
		return err
	}

	return nil
}

// ConfirmDeleteAPI 确认删除API
func ConfirmDeleteAPI(g *gocui.Gui, v *gocui.View) error {
	// 删除选中的API
	models.APIList = append(models.APIList[:models.SelectedAPI], models.APIList[models.SelectedAPI+1:]...)

	// 如果删除后列表为空，重置选中索引
	if len(models.APIList) == 0 {
		models.SelectedAPI = -1
	} else if models.SelectedAPI >= len(models.APIList) {
		// 如果删除的是最后一个API，选中前一个
		models.SelectedAPI = len(models.APIList) - 1
	}

	// 更新视图
	UpdateAPIList(g)

	// 清除确认提示
	statusView, _ := g.View("status")
	statusView.Clear()

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
