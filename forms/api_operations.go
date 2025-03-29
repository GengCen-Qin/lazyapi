package forms

import (
	"fmt"
	"strings"
	"unicode"

	"lazyapi/models/db"
	"lazyapi/models/entity"
	"lazyapi/models/service"
	"lazyapi/navigation"
	"lazyapi/utils"

	"github.com/GengCen-Qin/gocui"
)

// HandleAPI 保存API（新建或编辑）
func HandleAPI(g *gocui.Gui, v *gocui.View) error {
	if !FormInfo.Active {
		return nil
	}
	if FormInfo.IsDelete {
		return nil
	}

	// 收集表单数据
	var name, path, method, params string
	nameView, _ := g.View("form-name")
	methodView, _ := g.View("form-method")
	pathView, _ := g.View("form-path")
	paramsView, _ := g.View("form-params")

	if !FormInfo.IsFastApi {
		name = strings.TrimSpace(nameView.Buffer())
		method = strings.TrimSpace(methodView.Buffer())
	}
	path = strings.TrimSpace(pathView.Buffer())
	params = strings.TrimSpace(paramsView.Buffer())

	if err := validateAPIForm(g, name, path, method, params); err != nil {
		statusView, _ := g.View("status")
		statusView.Clear()
		fmt.Fprint(statusView, err.Error())
		return nil
	}

	if FormInfo.IsFastApi {
		tmpApi := entity.API{
			Path:   path,
			Params: params,
			Method: entity.MethodTitle[FormInfo.FastMethod],
		}
		param, _ := tmpApi.GetParams()
		sendRequest(g, &tmpApi, param)
		SelectLastRequestRecord(g)
		navigation.NavigateToRecordView(g)
		UpdateRequestRecordList(g)
	} else if FormInfo.IsEditing {
		service.EditAPI(entity.SelectedAPI, name, path, method, params)
	} else {
		newAPI := service.NewAPI(name, path, method, params)
		entity.SelectedAPI = newAPI.Id
	}

	// 更新视图
	UpdateAPIList(g)
	if FormInfo.IsFastApi {
		UpdateRequestRecordList(g)
	}

	// 关闭表单
	FormInfo.IsEditing = false
	// 关闭光标
	g.Cursor = false
	return CloseForm(g, v)
}

// EditAPIForm 编辑API表单
func EditAPIForm(g *gocui.Gui, v *gocui.View) error {
	if entity.SelectedAPI == -1 {
		return nil
	}

	// 标记为编辑模式
	FormInfo.IsEditing = true

	// 显示表单并填充数据
	if err := ShowNewAPIForm(g, v); err != nil {
		return err
	}

	// 获取当前选中的API
	api, _ := db.FindAPI(entity.SelectedAPI)

	fillFormFields(g, api)
	return nil
}

func FormatApiParam(g *gocui.Gui, v *gocui.View) error {
	view, err := g.View("form-params")
	if err != nil {
		return err
	}

	// 获取当前参数内容
	params := strings.TrimSpace(view.Buffer())
	if params == "" {
		return nil // 如果为空，不需要格式化
	}

	// 保存当前光标位置
	cx, cy := view.Cursor()
	
	// 计算当前光标在整个文本中的绝对位置
	absolutePos := calculateAbsolutePosition(params, cx, cy)
	
	// 尝试格式化JSON
	formattedJSON, err := utils.FormatJSON(params)
	if err != nil {
		// JSON格式有误，在状态栏显示错误信息
		statusView, _ := g.View("status")
		if statusView != nil {
			statusView.Clear()
			fmt.Fprintf(statusView, "JSON格式错误: %v", err)
		}
		return nil
	}

	// 清空视图并写入格式化后的JSON
	view.Clear()
	fmt.Fprint(view, formattedJSON)
	
	// 尝试将光标定位到格式化后文本中相对应的位置
	newCx, newCy := findBestCursorPosition(params, formattedJSON, absolutePos)
	view.SetCursor(newCx, newCy)

	return nil
}

// calculateAbsolutePosition 计算光标在整个文本中的绝对位置
func calculateAbsolutePosition(text string, x, y int) int {
	lines := strings.Split(text, "\n")
	position := 0
	
	// 计算前面几行的总字符数
	for i := 0; i < y && i < len(lines); i++ {
		position += len(lines[i]) + 1 // +1 是换行符
	}
	
	// 加上当前行的位置
	if y < len(lines) && x <= len(lines[y]) {
		position += x
	}
	
	return position
}

// findBestCursorPosition 在格式化后的文本中找到最佳光标位置
func findBestCursorPosition(originalText, formattedText string, originalPos int) (int, int) {
	// 计算原始文本中光标位置前的有效JSON内容（非空白字符）
	originalContentBeforeCursor := countNonWhitespaceChars(originalText, originalPos)
	
	// 在格式化后的文本中找到对应数量的非空白字符后的位置
	lines := strings.Split(formattedText, "\n")
	currentCount := 0
	
	for y, line := range lines {
		for x, char := range line {
			if !unicode.IsSpace(char) {
				currentCount++
			}
			
			if currentCount > originalContentBeforeCursor {
				// 找到了对应位置
				return x, y
			}
		}
	}
	
	// 如果没找到合适位置，返回文本末尾
	lastLineIndex := len(lines) - 1
	if lastLineIndex >= 0 {
		return len(lines[lastLineIndex]), lastLineIndex
	}
	
	return 0, 0
}

// countNonWhitespaceChars 计算文本中指定位置前的非空白字符数量
func countNonWhitespaceChars(text string, position int) int {
	count := 0
	for i, char := range text {
		if i >= position {
			break
		}
		
		if !unicode.IsSpace(char) {
			count++
		}
	}
	return count
}