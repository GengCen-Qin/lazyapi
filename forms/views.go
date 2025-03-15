package forms

import (
	"fmt"
	"strings"

	"lazyapi/common"

	"github.com/GengCen-Qin/gocui"
)

// SetCurrentViewOnTop 设置当前视图为顶层并更新状态栏
// 返回当前视图和可能的错误
func SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	// 检查视图是否存在
	if _, err := g.View(name); err != nil {
		return nil, fmt.Errorf("视图 '%s' 不存在: %v", name, err)
	}

	// 设置当前视图
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, fmt.Errorf("无法设置当前视图为 '%s': %v", name, err)
	}

	// 当视图变化时更新状态栏
	if err := UpdateStatusBar(g, name); err != nil {
		return nil, err
	}

	// 处理表单相关视图
	if isFormView(name) {
		if err := ensureFormOnTop(g); err != nil {
			return nil, err
		}
	}

	return g.SetViewOnTop(name)
}

// isFormView 检查是否为表单相关视图
func isFormView(name string) bool {
	return FormInfo.Active && strings.HasPrefix(name, "form")
}

// ensureFormOnTop 确保表单及其所有元素在顶部
func ensureFormOnTop(g *gocui.Gui) error {
	// 将表单容器置于顶部
	if _, err := g.SetViewOnTop("form"); err != nil {
		return err
	}

	// 将表单按钮置于顶部
	if _, err := g.SetViewOnTop("form-buttons"); err != nil {
		return err
	}

	// 将所有表单字段置于顶部
	for _, field := range FormInfo.Fields {
		fieldName := "form-" + field
		if _, err := g.SetViewOnTop(fieldName); err != nil {
			return err
		}
	}

	return nil
}

// NextView 切换到下一个视图
// 循环遍历所有主视图
func NextView(g *gocui.Gui, v *gocui.View) error {
	// 计算下一个视图索引
	nextIndex := (common.ViewActiveIndex + 1) % len(common.ViewArr)
	name := common.ViewArr[nextIndex]

	// 设置下一个视图为当前视图
	if _, err := SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	// 更新当前活动视图索引
	common.ViewActiveIndex = nextIndex

 	// 重新绘制界面以更新布局比例
    g.Update(func(g *gocui.Gui) error {
   		UpdateAPIList(g)
   		UpdateRequestRecordList(g)
        return nil // 触发重绘
    })
	return nil
}

func NextInfoView(g *gocui.Gui, v *gocui.View) error {
	// 计算下一个视图索引
	nextIndex := (common.InfoViewActiveIndex + 1) % len(common.InfoViewArr)
	name := common.InfoViewArr[nextIndex]

	// 设置下一个视图为当前视图
	if _, err := SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	// 更新当前活动视图索引
	common.InfoViewActiveIndex = nextIndex
	return nil
}

// UpdateStatusBar 更新状态栏信息
// 根据当前活动视图显示相应的提示信息
func UpdateStatusBar(g *gocui.Gui, viewName string) error {
	statusView, err := g.View("status")
	if err != nil {
		return fmt.Errorf("无法获取状态栏视图: %v", err)
	}

	statusView.Clear()

	// 获取并显示当前视图对应的状态消息
	message, exists := common.StatusMessages[viewName]
	if !exists {
		message = "Ctrl+C 退出"
	}

	fmt.Fprintf(statusView, "%s", message)
	return nil
}
