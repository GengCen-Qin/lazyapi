package utils

import (
	// "strings"

	"github.com/GengCen-Qin/gocui"
)

// 向上滚动视图
func ScrollViewUp(g *gocui.Gui, v *gocui.View) error {
	targetView, err := g.View("right-bottom")
	if err != nil {
		return err
	}
    scrollView(targetView, -1)
    return nil
}

// 向下滚动视图
func ScrollViewDown(g *gocui.Gui, v *gocui.View) error {
	targetView, err := g.View("right-bottom")
	if err != nil {
		return err
	}
    scrollView(targetView, 1)
    return nil
}

// 滚动视图的辅助函数
func scrollView(v *gocui.View, dy int) {
	if v != nil {
        ox, oy := v.Origin()
        if oy+dy >= 0 {  // 防止滚动到负位置
            v.SetOrigin(ox, oy+dy)
        }
    }
}

// 恢复视图位置
func ResetViewOrigin(v *gocui.View) {
    if v != nil {
        v.SetOrigin(0, 0)
    }
}
