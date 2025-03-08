package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"lazyapi/forms"
	"lazyapi/ui"
	"lazyapi/models"
)

func main() {
	defer models.CloseDB()
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true

	// 初始化UI
	if err := ui.Initialize(g); err != nil {
		log.Panicln(err)
	}

	// 设置表单键绑定
	if err := forms.SetupFormKeybindings(g); err != nil {
		log.Panicln(err)
	}

	// 使用Update方法确保在GUI完全初始化后更新API列表
	g.Update(func(g *gocui.Gui) error {
		forms.UpdateAPIList(g)
		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
