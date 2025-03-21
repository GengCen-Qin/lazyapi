package utils

import (
	"github.com/GengCen-Qin/gocui"
)

func ScrollViewPoint(g *gocui.Gui, v *gocui.View, view_name string, direact int) error {
	targetView, err := g.View(view_name)
	if err != nil {
		return err
	}
    return ScrollView(targetView, direact)
}

func ScrollApiInfoViewUp(g *gocui.Gui, v *gocui.View) error {
    return ScrollViewPoint(g, v, "api_info", -1)
}

func ScrollApiInfoViewDown(g *gocui.Gui, v *gocui.View) error {
	return ScrollViewPoint(g, v, "api_info", 1)
}

func ScrollRespondInfoViewUp(g *gocui.Gui, v *gocui.View) error {
	return ScrollViewPoint(g, v, "respond_info", -1)
}

func ScrollRespondInfoViewDown(g *gocui.Gui, v *gocui.View) error {
	return ScrollViewPoint(g, v, "respond_info", 1)
}

// ScrollView 通用滚动视图函数
// 按指定的增量调整视图的原点
func ScrollView(v *gocui.View, dy int) error {
	if v == nil {
        return nil
    }

    ox, oy := v.Origin()

    if oy+dy >= 0 {  // 防止滚动到负位置

        if err := v.SetOrigin(ox, oy+dy); err != nil {
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
