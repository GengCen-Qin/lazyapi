package common

// 视图相关的常量和状态
var (
	// List - 主要视图列表
	ViewArr = []string{"api_list", "record_list"}

	ViewIndexMap = map[string]int{
		"api_list":            0,
		"record_list": 1,
	}

	// Active - 当前激活的视图索引
	ViewActiveIndex = 0

	// StatusMessages - 每个视图对应的状态栏文案
	StatusMessages = map[string]string{
		"api_list":             "API_LIST | n(new), e(edit), d(delete), r(request), space(jump detail)",
		"api_info":        "REQUEST_DEFINITION",
		"respond_info":     "RESPOND_RESULT",
		"request_confirm_view": "ctrl-r(confirm), ctrl-q(cancel)",
		"record_list":  "d(delete)",
	}
)
