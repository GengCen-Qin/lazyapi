package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	// 视图分布 - 只有三个视图
	viewArr = []string{"left", "right-top", "right-bottom"}
	active  = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	out, err := g.View("right-top")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	// 所有视图都是可编辑的
	g.Cursor = true

	active = nextIndex
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// 左边视图占宽度的 3/10
	leftWidth := int(float64(maxX) * 0.3)

	// 左边视图: 占整个高度，宽度为 3/10
	if v, err := g.SetView("left", 0, 0, leftWidth-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "接口列表"
		v.Wrap = true

		if _, err = setCurrentViewOnTop(g, "left"); err != nil {
			return err
		}
	}

	// 右上视图: 宽度 7/10，高度为 1/2
	if v, err := g.SetView("right-top", leftWidth, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "接口定义"
		v.Wrap = true
		v.Autoscroll = true
	}

	// 右下视图: 宽度 7/10，高度为 1/2
	if v, err := g.SetView("right-bottom", leftWidth, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "响应定义"
		v.Editable = true
		v.Wrap = true
		fmt.Fprint(v, "Press TAB to change current view")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
