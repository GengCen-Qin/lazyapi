package common

// FormState 表单状态
type FormState struct {
	Active       bool // 激活模式
	CurrentField int
	Fields       []string
	Labels       map[string]string
	IsEditing    bool // 是否处于编辑模式
	IsDelete     bool // 是否处于删除模式
}

// 全局可访问的表单状态
var FormInfo = FormState{
	Active:       false,
	CurrentField: 0,
	Fields:       []string{"name", "method", "path", "params"},
	Labels: map[string]string{
		"name":   "API 名称",
		"path":   "请求路径",
		"method": "请求方式",
		"params": "请求参数(JSON)",
	},
	IsEditing: false, // 默认不是编辑模式
	IsDelete:  false, // 默认不是删除模式
}

// ViewInfo 视图信息
var (
	// ViewArr - 主要视图列表
	ViewArr = []string{"left", "request-history"}
	Active  = 0

	// StatusMessages - 每个视图对应的状态栏文案
	StatusMessages = map[string]string{
		"left":         "API_LIST | n(new), e(edit), d(delete), r(request)",
		"right-top":    "REQUEST_DEFINITION",
		"right-bottom": "RESPOND_RESULT",
		"requestConfirmView": "ctrl-r(confirm), ctrl-q(cancel)",
	}
)
