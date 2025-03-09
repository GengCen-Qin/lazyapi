package ui

import (
	"fmt"

	"github.com/GengCen-Qin/gocui"
	"github.com/atotto/clipboard"
)

// Layout 管理GUI布局
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// 底部状态栏: 横跨整个宽度
	if v, err := g.SetView("status", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Frame = false  // 不显示边框
		v.Editable = false
		v.FgColor = gocui.ColorBlue // 设置前景色
	}

	// 左边视图占宽度的 1/5 (20%)
	leftWidth := int(float64(maxX) * 0.2)

	// 左边视图: 占整个高度，宽度为 1/5
	if v, err := g.SetView("left", 0, 0, leftWidth-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "接口列表"
		v.Wrap = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Editable = false

		if _, err = SetCurrentViewOnTop(g, "left"); err != nil {
			return err
		}
	}

	// 右上视图: 宽度 4/5，高度为 1/2
	if v, err := g.SetView("right-top", leftWidth, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "接口定义"
		v.Wrap = true
		v.Autoscroll = true
		v.Editable = false
	}

	// 右下视图: 宽度 4/5，高度为 1/2
	if v, err := g.SetView("right-bottom", leftWidth, maxY/2, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "响应展示"
		v.Editable = false
		v.Wrap = true
	}

	if err := g.SetKeybinding("right-bottom", 'y', gocui.ModNone, copyResponseToClipboard); err != nil {
		return err
	}

	return nil
}


func copyResponseToClipboard(g *gocui.Gui, v *gocui.View) error {
    if v == nil {
        return nil
    }

    // 获取right-bottom视图
    responseView, err := g.View("right-bottom")
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
