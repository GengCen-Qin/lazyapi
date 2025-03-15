package forms

import (
	"fmt"
	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/models/service"

	"github.com/GengCen-Qin/gocui"
)

// DeleteConfirmation 删除确认的结构体
type DeleteConfirmation struct {
	EntityType    string                            // 实体类型 ("API", "RequestRecord")
	EntityID      int                               // 要删除的实体ID
	DeleteFunc    func(int) error                   // 执行删除的函数
	AfterDelete   func(g *gocui.Gui) error          // 删除后的回调函数
	SuccessMsg    string                            // 删除成功的消息
}

var currentDeleteConfirmation *DeleteConfirmation

// ShowDeleteConfirmation 显示通用的删除确认对话框
func ShowDeleteConfirmation(g *gocui.Gui, confirmation *DeleteConfirmation) error {
	// 设置删除标记
	FormInfo.IsDelete = true

	// 保存当前删除确认信息
	currentDeleteConfirmation = confirmation

	// 显示确认提示
	statusView, _ := g.View("status")
	statusView.Clear()
	fmt.Fprintf(statusView, "Confirm to delete %s? (y/n)", confirmation.EntityType)

	// 绑定确认和取消操作
	if err := g.SetKeybinding("", 'y', gocui.ModNone, ConfirmDelete); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'n', gocui.ModNone, CancelDelete); err != nil {
		return err
	}

	return nil
}

// ConfirmDelete 通用的确认删除回调
func ConfirmDelete(g *gocui.Gui, v *gocui.View) error {
	if currentDeleteConfirmation == nil {
		return nil
	}

	// 执行删除操作
	if err := currentDeleteConfirmation.DeleteFunc(currentDeleteConfirmation.EntityID); err != nil {
		statusView, _ := g.View("status")
		statusView.Clear()
		fmt.Fprintf(statusView, "Delete failed: %v", err)
		cleanupDeleteBindings(g)
		return nil
	}

	// 执行删除后的回调
	if currentDeleteConfirmation.AfterDelete != nil {
		if err := currentDeleteConfirmation.AfterDelete(g); err != nil {
			return err
		}
	}

	// 显示成功消息
	statusView, _ := g.View("status")
	statusView.Clear()
	fmt.Fprint(statusView, currentDeleteConfirmation.SuccessMsg)

	// 清理
	cleanupDeleteBindings(g)
	return nil
}

// CancelDelete 通用的取消删除回调
func CancelDelete(g *gocui.Gui, v *gocui.View) error {
	// 清除确认提示
	statusView, _ := g.View("status")
	statusView.Clear()

	// 清理
	cleanupDeleteBindings(g)
	return nil
}

// cleanupDeleteBindings 清理删除操作的临时状态
func cleanupDeleteBindings(g *gocui.Gui) {
	// 删除临时键绑定
	g.DeleteKeybinding("", 'y', gocui.ModNone)
	g.DeleteKeybinding("", 'n', gocui.ModNone)

	// 重置删除标志
	FormInfo.IsDelete = false
	currentDeleteConfirmation = nil
}

// DeleteAPI 删除选中的API (使用新的通用确认机制)
func DeleteAPI(g *gocui.Gui, v *gocui.View) error {
	if entity.SelectedAPI == -1 {
		return nil
	}

	confirmation := &DeleteConfirmation{
		EntityType: "API",
		EntityID:   entity.SelectedAPI,
		DeleteFunc: db.DeleteApi,
		AfterDelete: func(g *gocui.Gui) error {
			list := service.APIList()
			if len(list) == 0 {
				entity.SelectedAPI = -1
			} else {
				entity.SelectedAPI = list[len(list)-1].Id
			}
			// 更新视图
			UpdateAPIList(g)
			return nil
		},
		SuccessMsg: "API deleted successfully!",
	}

	return ShowDeleteConfirmation(g, confirmation)
}

// DeleteRequestRecord 删除选中的请求记录 (使用通用确认机制)
func DeleteRequestRecord(g *gocui.Gui, v *gocui.View) error {
	if entity.SelectedQuestRecord == -1 {
		return nil
	}

	confirmation := &DeleteConfirmation{
		EntityType: "Request Record",
		EntityID:   entity.SelectedQuestRecord,
		DeleteFunc: db.DeleteRecord,
		AfterDelete: func(g *gocui.Gui) error {
			list := service.RequestRecordList()
			if len(list) == 0 {
				entity.SelectedQuestRecord = -1
			} else {
				entity.SelectedQuestRecord = list[len(list)-1].Id
			}
			// 更新请求记录列表视图
			UpdateRequestRecordList(g) // 假设这个函数存在
			return nil
		},
		SuccessMsg: "Request record deleted successfully!",
	}

	return ShowDeleteConfirmation(g, confirmation)
}
