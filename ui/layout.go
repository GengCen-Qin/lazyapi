package ui

import (
	"github.com/GengCen-Qin/gocui"
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

	return nil
}
