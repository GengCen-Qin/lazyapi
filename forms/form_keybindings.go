package forms

import (
	"github.com/GengCen-Qin/gocui"
)

// SetupFormKeybindings 为表单设置键绑定
func SetupFormKeybindings(g *gocui.Gui) error {
	// 左侧视图键绑定 - 'n'键创建新API或取消删除
	if err := g.SetKeybinding("api_list", 'n', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if FormInfo.IsDelete {
			return CancelDelete(g, v)
		}
		return ShowNewAPIForm(g, v)
	}); err != nil {
		return err
	}

	// 左侧视图键绑定 - 'e'键编辑选中的API
	if err := g.SetKeybinding("api_list", 'e', gocui.ModNone, EditAPIForm); err != nil {
		return err
	}

	// 左侧视图键绑定 - 'd'键删除选中的API
	if err := g.SetKeybinding("api_list", 'd', gocui.ModNone, DeleteAPI); err != nil {
		return err
	}

	// 右上视图键绑定 - 'r'键请求当前API
	if err := g.SetKeybinding("api_list", 'r', gocui.ModNone, RequestAPI); err != nil {
		return err
	}

	if err := g.SetKeybinding("api_list", gocui.KeySpace, gocui.ModNone, JumpApiDetail); err != nil {
		return err
	}

	if err := g.SetKeybinding("record_list", gocui.KeySpace, gocui.ModNone, JumpApiDetail); err != nil {
		return err
	}

	// 添加上下键绑定
	if err := g.SetKeybinding("api_list", gocui.KeyArrowUp, gocui.ModNone, MoveAPISelectionUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("api_list", gocui.KeyArrowDown, gocui.ModNone, MoveAPISelectionDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("api_info", gocui.KeyEsc, gocui.ModNone, JumpOut); err != nil {
		return err
	}

	if err := g.SetKeybinding("respond_info", gocui.KeyEsc, gocui.ModNone, JumpOut); err != nil {
		return err
	}

	if err := g.SetKeybinding("record_list", gocui.KeyArrowUp, gocui.ModNone, MoveRequestRecordSelectionUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("record_list", gocui.KeyArrowDown, gocui.ModNone, MoveRequestRecordSelectionDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("record_list", 'd', gocui.ModNone, DeleteRequestRecord); err != nil {
		return err
	}

	return nil
}
