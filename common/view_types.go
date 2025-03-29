package common

// 视图相关的常量和状态
var (
	// 列表视图
	ViewArr = []string{"api_list", "record_list"}

	ViewIndexMap = map[string]int{
		"api_list": 0,
		"record_list": 1,
	}

	ViewActiveIndex = 0

	// 详情视图
	InfoViewArr = []string{"api_info", "respond_info"}

	InfoViewIndexMap = map[string]int{
		"api_info": 0,
		"respond_info": 1,
	}

	InfoViewActiveIndex = 0

	// StatusMessages - 每个视图对应的状态栏文案
	StatusMessages = map[string]string{
		"api_list":             "API_LIST | n(new), e(edit), d(delete), r(request), space(jump detail)，tab(switch view), g(fast get request), p(fast post request)",
		"api_info":        "API_INFO | ↑(page up), ↓(page down), esc(back list)",
		"respond_info":     "RESPOND_INFO | ↑(page up), ↓(page down), esc(back list)",
		"request_confirm_view": "ctrl-r(confirm), ctrl-q(cancel)",
		"record_list":  "RECORD_LIST | d(delete), space(jump detail), tab(switch view), g(fast get request), p(fast post request)",
		"form-params": "ctrl-f(format), esc(concel)",
		"form-name": "esc(concel)",
		"form-method": "esc(concel)",
		"form-path": "esc(concel)",
	}
)
