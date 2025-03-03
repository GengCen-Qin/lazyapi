package common

// FormState 表单状态
type FormState struct {
	Active       bool
	CurrentField int
	Fields       []string
	Labels       map[string]string
}

// 全局可访问的表单状态
var FormInfo = FormState{
	Active:       false,
	CurrentField: 0,
	Fields:       []string{"name", "path", "method"},
	Labels: map[string]string{
		"name":   "API 名称",
		"path":   "请求路径",
		"method": "请求方式",
	},
}

// ViewInfo 视图信息
var (
	// ViewArr - 主要视图列表
	ViewArr = []string{"left", "right-top", "right-bottom"}
	Active  = 0

	// StatusMessages - 每个视图对应的状态栏文案
	StatusMessages = map[string]string{
		"left":         "接口列表 | n(new), e(edit), d(delete)",
		"right-top":    "接口定义 | s(save), c(cancel)",
		"right-bottom": "响应定义 | r(request), f(format)",
	}
)
