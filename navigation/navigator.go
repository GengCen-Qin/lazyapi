package navigation

import (
	"lazyapi/common"
	
	"github.com/GengCen-Qin/gocui"
)

// NavigateToView 切换到指定视图
func NavigateToView(g *gocui.Gui, viewName string) error {
	// 根据视图名称切换到对应视图
	if _, ok := common.ViewIndexMap[viewName]; ok {
		common.ViewActiveIndex = common.ViewIndexMap[viewName]
		return activateView(g, viewName)
	}
	
	if _, ok := common.InfoViewIndexMap[viewName]; ok {
		common.InfoViewActiveIndex = common.InfoViewIndexMap[viewName]
		return activateView(g, viewName)
	}
	
	return nil
}

// NavigateToRecordView 导航到记录视图
func NavigateToRecordView(g *gocui.Gui) error {
	common.ViewActiveIndex = common.ViewIndexMap["record_list"]
	return activateView(g, "record_list")
}

// NavigateToRespondInfo 导航到响应信息视图
func NavigateToRespondInfo(g *gocui.Gui) error {
	common.InfoViewActiveIndex = common.InfoViewIndexMap["respond_info"]
	return activateView(g, "respond_info")
}

// 激活指定视图
func activateView(g *gocui.Gui, name string) error {
	if v, err := g.View(name); err == nil {
		if _, err := g.SetCurrentView(name); err != nil {
			return err
		}
		v.Highlight = true
		
		// 更新状态栏
		statusView, err := g.View("status")
		if err != nil {
			return err
		}
		statusView.Clear()
		if msg, ok := common.StatusMessages[name]; ok {
			statusView.Write([]byte(msg))
		}
	}
	g.Update(func(g *gocui.Gui) error {
		return nil
	})

	return nil
}
