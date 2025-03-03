package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	// 初始GUI
	g := initializeGUI()
	// 函数结束后关闭
	defer g.Close()
	// 设置UI布局
	setupGUI(g)
	// 一直监听直到有异常发生
	runMainLoop(g)
}


func initializeGUI() *gocui.Gui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	return g
}

func setupGUI(g *gocui.Gui) {
	g.SetManagerFunc(layout)

	// 设置快捷键退出
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func runMainLoop(g *gocui.Gui) {
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
