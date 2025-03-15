package ui

import (
	"fmt"
	"lazyapi/common"
	"lazyapi/forms"

	"github.com/GengCen-Qin/gocui"
	"github.com/atotto/clipboard"
)

// Layout 管理GUI布局
// 负责创建并定位所有视图
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// 创建各个主要区域
	if err := createStatusBar(g, maxX, maxY); err != nil {
		return err
	}

	// 左侧视图占宽度的 1/3 (33.3%)
	leftWidth := int(float64(maxX) * 0.333)

	// 根据当前活动视图决定左侧各视图的高度比例
	var leftApiHeight int
	currentView := common.ViewArr[common.ViewActiveIndex]

	// 计算左侧视图的高度分配
	if currentView == "api_list" {
		// 如果当前是API列表视图，则它占70%
		leftApiHeight = int(float64(maxY-2) * 0.7)
	} else if currentView == "record_list" {
		// 如果当前是请求历史视图，则它占70%
		leftApiHeight = int(float64(maxY-2) * 0.3)
	} else {
		// 默认情况下，API列表占70%
		leftApiHeight = int(float64(maxY-2) * 0.7)
	}

	if err := createApiListView(g, leftWidth, leftApiHeight); err != nil {
		return err
	}

	if err := createHistoryView(g, leftWidth, leftApiHeight, maxY); err != nil {
		return err
	}

	if err := createDetailViews(g, leftWidth, maxX, maxY); err != nil {
		return err
	}

	return nil
}

// createStatusBar 创建底部状态栏
func createStatusBar(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("status", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Frame = false  // 不显示边框
		v.Editable = false
		v.FgColor = gocui.ColorBlue // 设置前景色
	}
	return nil
}

// createApiListView 创建API列表视图（左上）
func createApiListView(g *gocui.Gui, leftWidth, leftApiHeight int) error {
	if v, err := g.SetView("api_list", 0, 0, leftWidth-1, leftApiHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "接口列表"
		v.Wrap = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Editable = false

		// 只有在初始化时才设置当前视图
		if common.ViewActiveIndex == 0 {
			if _, err = forms.SetCurrentViewOnTop(g, "api_list"); err != nil {
				return err
			}
		}
	}
	return nil
}

// createHistoryView 创建请求历史视图（左下）
func createHistoryView(g *gocui.Gui, leftWidth, leftApiHeight, maxY int) error {
	if v, err := g.SetView("record_list", 0, leftApiHeight+1, leftWidth-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "请求记录"
		v.Editable = false
	}
	return nil
}

// createDetailViews 创建详情视图（右侧）
func createDetailViews(g *gocui.Gui, leftWidth, maxX, maxY int) error {
	// 右上视图: 接口定义
	if v, err := g.SetView("api_info", leftWidth, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "接口定义"
		v.Wrap = true
		v.Autoscroll = false
		v.Editable = false
	}

	// 右下视图: 响应展示
	if v, err := g.SetView("respond_info", leftWidth, maxY/2, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "响应展示"
		v.Editable = false
		v.Wrap = true
		v.Autoscroll = false
	}

	return nil
}

// CopyResponseToClipboard 复制响应内容到剪贴板
func CopyResponseToClipboard(g *gocui.Gui, v *gocui.View) error {
    if v == nil {
        return nil
    }

    // 获取respond_info视图
    responseView, err := g.View("respond_info")
    if err != nil {
        return err
    }

    // 获取视图中的所有内容
    responseContent := responseView.Buffer()
    // 复制到剪贴板
    err = clipboard.WriteAll(responseContent)
    if err != nil {
        return err
    }

	statusView, _ := g.View("status")
    statusView.Clear()
    fmt.Fprint(statusView, "响应内容已复制到剪贴板")

    return nil
}
