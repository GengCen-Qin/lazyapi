package common

// 视图相关的常量和状态
var (
	// List - 主要视图列表
	ViewArr = []string{"left", "request-history"}

	ViewIndexMap = map[string]int{
		"left":            0,
		"request-history": 1,
	}

	// Active - 当前激活的视图索引
	ViewActiveIndex = 0

	// StatusMessages - 每个视图对应的状态栏文案
	StatusMessages = map[string]string{
		"left":             "API_LIST | n(new), e(edit), d(delete), r(request)",
		"right-top":        "REQUEST_DEFINITION",
		"right-bottom":     "RESPOND_RESULT",
		"requestConfirmView": "ctrl-r(confirm), ctrl-q(cancel)",
		"request-history":  "d(delete)",
	}
)
