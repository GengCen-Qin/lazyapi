package utils

import (
	"fmt"
	"strings"
	"unicode"
)

// TruncateString 截断字符串到指定长度，可选添加省略号
// 常用于确保字符串在有限宽度的视图中正确显示
func TruncateString(s string, maxLength int, addEllipsis bool) string {
    if len(s) <= maxLength {
        return s
    }

    if addEllipsis && maxLength > 3 {
        return s[:maxLength-3] + "..."
    }

    return s[:maxLength]
}

// PadRight 右侧填充字符串到指定长度
// 常用于表格式布局中保持列宽一致
func PadRight(s string, padChar rune, length int) string {
    if len(s) >= length {
        return s
    }

    return s + strings.Repeat(string(padChar), length-len(s))
}

// PadLeft 左侧填充字符串到指定长度
func PadLeft(s string, padChar rune, length int) string {
    if len(s) >= length {
        return s
    }

    return strings.Repeat(string(padChar), length-len(s)) + s
}

// RemoveWhitespace 移除字符串中的所有空白字符
func RemoveWhitespace(s string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            return -1
        }
        return r
    }, s)
}

// HighlightText 突出显示文本(用ANSI颜色代码)
// 可指定颜色代码：31(红), 32(绿), 33(黄), 34(蓝), 35(紫), 36(青), 37(白)
func HighlightText(text string, colorCode int) string {
    if colorCode < 30 || colorCode > 37 {
        colorCode = 31 // 默认红色
    }
    return fmt.Sprintf("\033[%d;1m%s\033[0m", colorCode, text)
}

// WrapInBox 将文本用边框包围
// 常用于在终端中突出显示某些信息
func WrapInBox(text string) string {
    lines := strings.Split(text, "\n")
    maxWidth := 0

    // 找出最长行的长度
    for _, line := range lines {
        if len(line) > maxWidth {
            maxWidth = len(line)
        }
    }

    // 上边框
    result := "+" + strings.Repeat("-", maxWidth+2) + "+\n"

    // 内容
    for _, line := range lines {
        result += "| " + PadRight(line, ' ', maxWidth) + " |\n"
    }

    // 下边框
    result += "+" + strings.Repeat("-", maxWidth+2) + "+"

    return result
}
