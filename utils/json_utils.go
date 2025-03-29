package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatJSON 格式化JSON字符串
// 将紧凑的JSON格式化为带缩进的易读格式，保持原始文本结构
func FormatJSON(jsonString string) (string, error) {
    if jsonString == "" {
        return "", nil
    }

    // 首先验证JSON是否有效
    if !json.Valid([]byte(jsonString)) {
        return jsonString, fmt.Errorf("无效的JSON格式")
    }

    // 使用字符级别的格式化，保持原始文本结构
    return formatJSONPreserveOrder(jsonString)
}

// formatJSONPreserveOrder 格式化JSON字符串，保持原始顺序和结构
func formatJSONPreserveOrder(jsonStr string) (string, error) {
    var result strings.Builder
    var indentLevel int
    var inString bool
    var escapeNext bool

    // 移除所有现有的空白字符
    jsonStr = removeWhitespace(jsonStr)

    for i, char := range jsonStr {
        if escapeNext {
            result.WriteRune(char)
            escapeNext = false
            continue
        }

        switch char {
        case '\\':
            result.WriteRune(char)
            escapeNext = true
        case '"':
            result.WriteRune(char)
            // 如果前一个字符不是转义字符，则切换字符串状态
            if i == 0 || jsonStr[i-1] != '\\' {
                inString = !inString
            }
        case '{', '[':
            result.WriteRune(char)
            if !inString {
                indentLevel++
                result.WriteString("\n")
                result.WriteString(strings.Repeat("    ", indentLevel))
            }
        case '}', ']':
            if !inString {
                indentLevel--
                result.WriteString("\n")
                result.WriteString(strings.Repeat("    ", indentLevel))
            }
            result.WriteRune(char)
        case ',':
            result.WriteRune(char)
            if !inString {
                result.WriteString("\n")
                result.WriteString(strings.Repeat("    ", indentLevel))
            }
        case ':':
            result.WriteRune(char)
            if !inString {
                result.WriteString(" ")
            }
        default:
            result.WriteRune(char)
        }
    }

    return result.String(), nil
}

// removeWhitespace 移除JSON字符串中的所有空白字符，但保留字符串内的空白
func removeWhitespace(jsonStr string) string {
    var result strings.Builder
    var inString bool
    var escapeNext bool

    for i, char := range jsonStr {
        if escapeNext {
            result.WriteRune(char)
            escapeNext = false
            continue
        }

        switch char {
        case '\\':
            result.WriteRune(char)
            escapeNext = true
        case '"':
            result.WriteRune(char)
            // 如果前一个字符不是转义字符，则切换字符串状态
            if i == 0 || jsonStr[i-1] != '\\' {
                inString = !inString
            }
        case ' ', '\t', '\n', '\r':
            // 只保留字符串内的空白
            if inString {
                result.WriteRune(char)
            }
        default:
            result.WriteRune(char)
        }
    }

    return result.String()
}

// IsValidJSON 检查字符串是否为有效的JSON
func IsValidJSON(str string) bool {
    return json.Valid([]byte(str))
}

// PrettyPrintJSON 格式化并输出JSON，适用于终端显示
// 增加颜色高亮等特性
func PrettyPrintJSON(jsonStr string) (string, error) {
    if jsonStr == "" {
        return "", nil
    }

    // 先格式化JSON字符串
    formatted, err := FormatJSON(jsonStr)
    if err != nil {
        return jsonStr, err
    }

    // 应用颜色高亮
    colorized, err := ColorizeJSON(formatted)
    if err != nil {
        // 如果高亮失败，至少返回格式化后的JSON
        return formatted, fmt.Errorf("格式化成功但高亮失败: %v", err)
    }

    return colorized, nil
}

// CompactJSON 压缩JSON字符串，移除所有空白字符
func CompactJSON(jsonStr string) (string, error) {
    if jsonStr == "" {
        return "", nil
    }

    // 验证JSON是否有效
    if !json.Valid([]byte(jsonStr)) {
        return jsonStr, fmt.Errorf("无效的JSON格式")
    }

    return removeWhitespace(jsonStr), nil
}
