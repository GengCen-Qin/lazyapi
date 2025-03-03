package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"lazyapi/forms"
	"lazyapi/ui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	// 初始化UI
	if err := ui.Initialize(g); err != nil {
		log.Panicln(err)
	}

	// 设置表单键绑定
	if err := forms.SetupFormKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
