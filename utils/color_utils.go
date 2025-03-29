package utils

import (
	"encoding/json"
	"github.com/hokaccha/go-prettyjson"
)

// 定义常用ANSI颜色代码
const (
    ColorReset  = "\033[0m"
    ColorBold   = "\033[1m"
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorPurple = "\033[35m"
    ColorCyan   = "\033[36m"
    ColorWhite  = "\033[37m"
    ColorMagenta = "\033[35m"

    BoldRed    = "\033[31;1m"
    BoldGreen  = "\033[32;1m"
    BoldYellow = "\033[33;1m"
    BoldBlue   = "\033[34;1m"
    BoldPurple = "\033[35;1m"
    BoldCyan   = "\033[36;1m"
    BoldWhite  = "\033[37;1m"
)

// ColorText 用指定颜色包装文本
func ColorText(text string, colorCode string) string {
    return colorCode + text + ColorReset
}

// RedText 红色文本
func RedText(text string) string {
    return ColorText(text, ColorRed)
}

// GreenText 绿色文本
func GreenText(text string) string {
    return ColorText(text, ColorGreen)
}

// YellowText 黄色文本
func YellowText(text string) string {
    return ColorText(text, ColorYellow)
}

// BlueText 蓝色文本
func BlueText(text string) string {
    return ColorText(text, ColorBlue)
}

// SuccessText 成功文本 (绿色加粗)
func SuccessText(text string) string {
    return ColorText(text, BoldGreen)
}

// ErrorText 错误文本 (红色加粗)
func ErrorText(text string) string {
    return ColorText(text, BoldRed)
}

// WarningText 警告文本 (黄色加粗)
func WarningText(text string) string {
    return ColorText(text, BoldYellow)
}

// InfoText 信息文本 (蓝色加粗)
func InfoText(text string) string {
    return ColorText(text, BoldBlue)
}

// ColorizeJSON 为JSON字符串添加语法高亮
func ColorizeJSON(jsonStr string) (string, error) {
    var data interface{}
    if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
        return "", err
    }

    formatter := prettyjson.NewFormatter()
    formatter.Indent = 2
    formatter.Newline = "\n"
    result, err := formatter.Marshal(data)

    return string(result), err
}
