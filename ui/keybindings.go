package ui

import (
	// "strings"

	"github.com/GengCen-Qin/gocui"
	"lazyapi/utils"
)

// Quit 退出程序
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// GlobleSetupKeybindings 设置键盘绑定
func GlobleSetupKeybindings(g *gocui.Gui) error {
	// 全局键绑定
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("left", gocui.KeyTab, gocui.ModNone, NextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("request-history", gocui.KeyTab, gocui.ModNone, NextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("left", 'y', gocui.ModNone, copyResponseToClipboard); err != nil {
		return err
	}
    if err := g.SetKeybinding("left", gocui.KeyCtrlLsqBracket, gocui.ModNone, utils.ScrollViewUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("left", gocui.KeyCtrlRsqBracket, gocui.ModNone, utils.ScrollViewDown); err != nil {
        return err
    }

	return nil
}

// Initialize 初始化UI
func Initialize(g *gocui.Gui) error {
	// 设置布局管理器
	g.SetManagerFunc(Layout)

	// 设置键盘绑定
	if err := GlobleSetupKeybindings(g); err != nil {
		return err
	}

	return nil
}
