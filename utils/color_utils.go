package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
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
    // 先验证JSON是否有效
    var tmp interface{}
    if err := json.Unmarshal([]byte(jsonStr), &tmp); err != nil {
        return "", fmt.Errorf("invalid JSON: %w", err)
    }

    var result strings.Builder
    var inString bool
    var escaped bool

    // 更精确的状态处理
    for i, c := range jsonStr {
        switch {
        case escaped:
            // 处理转义字符后的字符
            result.WriteRune(c)
            escaped = false

        case c == '\\':
            // 标记转义序列的开始
            result.WriteRune(c)
            escaped = true

        case c == '"':
            // 处理字符串边界
            if inString {
                result.WriteString(string(c) + ColorReset)
                inString = false
            } else {
                result.WriteString(ColorRed + string(c))
                inString = true
            }

        case c == ':' && !inString:
            // 处理键值分隔符
            result.WriteString(ColorReset + string(c) + " ")

            // 预先查看下一个非空白字符，为值选择合适的颜色
            nextColor := ColorBlue // 默认为蓝色
            for j := i + 1; j < len(jsonStr); j++ {
                next := rune(jsonStr[j])
                if unicode.IsSpace(next) {
                    continue
                }

                switch next {
                case '"': // 字符串值将在下一个循环中处理
                    nextColor = ""
                case '{', '[': // 对象或数组不需要特殊颜色
                    nextColor = ColorReset
                case 't', 'f', 'n': // true, false, null
                    nextColor = ColorGreen
                default: // 数字或其他
                    if next >= '0' && next <= '9' || next == '-' || next == '+' {
                        nextColor = ColorMagenta
                    }
                }
                break
            }

            if nextColor != "" {
                result.WriteString(nextColor)
            }

        case (c == '{' || c == '}' || c == '[' || c == ']') && !inString:
            // 处理结构符号，添加颜色重置和格式化空格
            if i > 0 && !unicode.IsSpace(rune(jsonStr[i-1])) && jsonStr[i-1] != ':' {
                result.WriteString(ColorReset + " ")
            }
            result.WriteString(ColorYellow + string(c) + ColorReset)

            // 在关闭括号后添加换行符以提高可读性（可选）
            if c == '}' || c == ']' {
                // 检查下一个字符是否为逗号或结束括号
                if i+1 < len(jsonStr) && (jsonStr[i+1] == ',' || jsonStr[i+1] == '}' || jsonStr[i+1] == ']') {
                    // 不添加换行
                } else {
                    result.WriteString("\n")
                }
            }

        case c == ',' && !inString:
            // 处理分隔符
            result.WriteString(ColorReset + string(c) + " ")

        default:
            // 处理常规字符
            if !inString && (c >= '0' && c <= '9' || c == '-' || c == '+' || c == '.') {
                // 数字用洋红色
                if i == 0 || !unicode.IsDigit(rune(jsonStr[i-1])) && jsonStr[i-1] != '.' && jsonStr[i-1] != '-' {
                    result.WriteString(ColorMagenta)
                }
                result.WriteRune(c)
                // 检查数字结束
                if i+1 < len(jsonStr) && !unicode.IsDigit(rune(jsonStr[i+1])) && jsonStr[i+1] != '.' {
                    result.WriteString(ColorReset)
                }
            } else if !inString && (strings.HasPrefix(jsonStr[i:], "true") ||
                                    strings.HasPrefix(jsonStr[i:], "false") ||
                                    strings.HasPrefix(jsonStr[i:], "null")) {
                // 布尔值和null用绿色
                var keyword string
                if strings.HasPrefix(jsonStr[i:], "true") {
                    keyword = "true"
                } else if strings.HasPrefix(jsonStr[i:], "false") {
                    keyword = "false"
                } else {
                    keyword = "null"
                }

                if i == 0 || !unicode.IsLetter(rune(jsonStr[i-1])) {
                    result.WriteString(ColorGreen + keyword + ColorReset)
                    i += len(keyword) - 1  // 跳过关键字的其余部分
                } else {
                    result.WriteRune(c)
                }
            } else {
                result.WriteRune(c)
            }
        }
    }

    // 确保结束时重置颜色
    result.WriteString(ColorReset)
    return result.String(), nil
}
