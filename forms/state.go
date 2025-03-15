package forms

// State 表单状态
type State struct {
	Active       bool            // 激活模式
	CurrentField int             // 当前选中的字段索引
	Fields       []string        // 表单字段列表
	Labels       map[string]string // 字段对应的显示标签
	IsEditing    bool            // 是否处于编辑模式
	IsDelete     bool            // 是否处于删除模式
}

// 创建新的默认表单状态
func NewDefaultState() State {
	return State{
		Active:       false,
		CurrentField: 0,
		Fields:       []string{"name", "method", "path", "params"},
		Labels: map[string]string{
			"name":   "API 名称",
			"path":   "请求路径",
			"method": "请求方式",
			"params": "请求参数(JSON)",
		},
		IsEditing: false,
		IsDelete:  false,
	}
}

// 全局可访问的表单状态
var FormInfo = NewDefaultState()
