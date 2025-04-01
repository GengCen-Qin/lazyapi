package utils

import (
	"strings"

	"github.com/GengCen-Qin/gocui"
)

const (
    ScrollUp   = -1
    ScrollDown = 1
)

var (
    ScrollApiInfoViewUp     = ScrollViewFactory("api_info", ScrollUp)
    ScrollApiInfoViewDown   = ScrollViewFactory("api_info", ScrollDown)
    ScrollRespondInfoViewUp   = ScrollViewFactory("respond_info", ScrollUp)
    ScrollRespondInfoViewDown = ScrollViewFactory("respond_info", ScrollDown)
)

func ScrollViewPoint(g *gocui.Gui, v *gocui.View, view_name string, direact int) error {
	targetView, err := g.View(view_name)
	if err != nil {
		return err
	}
    return ScrollView(targetView, direact)
}

func ScrollViewFactory(viewName string, direction int) func(*gocui.Gui, *gocui.View) error {
    return func(g *gocui.Gui, v *gocui.View) error {
        return ScrollViewPoint(g, v, viewName, direction)
    }
}

// ScrollView 通用滚动视图函数
// 按指定的增量调整视图的原点
func ScrollView(v *gocui.View, dy int) error {
    if v == nil {
        return nil
    }

    ox, oy := v.Origin()
    _, viewHeight := v.Size()

    // 获取视图内容的总行数
    contentHeight := len(strings.Split(v.Buffer(), "\n")) - 1 // -1 因为最后一个换行符会产生一个空行

    // 计算最大可滚动位置(如果内容不足以填满视图，则不需要滚动)
    maxScroll := contentHeight - viewHeight
    if maxScroll < 0 {
        maxScroll = 0
    }

    // 计算新的原点位置，确保不超出有效范围
    newOy := oy + dy
    if newOy < 0 {
        newOy = 0 // 防止向上滚动超出顶部
    } else if newOy > maxScroll {
        newOy = maxScroll // 防止向下滚动超出底部
    }

    // 只有当位置确实需要改变时才设置新的原点
    if newOy != oy {
        if err := v.SetOrigin(ox, newOy); err != nil {
            return err
        }
    }

    return nil
}

// ResetViewOrigin 重置视图的原点位置
// 通常在清空视图内容后使用
func ResetViewOrigin(v *gocui.View) error {
    if v == nil {
        return nil
    }
    return v.SetOrigin(0, 0)
}

// ClearView 清空视图内容并重置原点
func ClearView(v *gocui.View) error {
    if v == nil {
        return nil
    }
    v.Clear()
    return ResetViewOrigin(v)
}

// ViewDimensions 获取视图的尺寸信息
// 返回视图的宽度和高度
func ViewDimensions(v *gocui.View) (width, height int) {
    if v == nil {
        return 0, 0
    }

    w, h := v.Size()
    return w, h
}

// IsViewEmpty 检查视图是否为空
func IsViewEmpty(v *gocui.View) bool {
    if v == nil {
        return true
    }

    return len(v.Buffer()) == 0
}
