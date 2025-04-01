package main

import (
	"log"

	"github.com/GengCen-Qin/gocui"
	"lazyapi/forms"
	"lazyapi/ui"
	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/models/service"
)

func main() {
	defer db.CloseDB()
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true
	g.Mouse = true

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
		apiList := service.APIList()
		if len(apiList) != 0 {
			entity.SelectedAPI = apiList[0].Id
		}
		recordList := service.RequestRecordList()
		if len(recordList) != 0 {
			entity.SelectedQuestRecord = recordList[0].Id
		}
 		forms.UpdateAPIList(g)
   		forms.UpdateRequestRecordList(g)
		return nil
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
