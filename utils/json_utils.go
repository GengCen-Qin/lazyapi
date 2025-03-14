package utils

import (
	"encoding/json"
	"fmt"
)

// FormatJSON 格式化JSON字符串
// 将紧凑的JSON格式化为带缩进的易读格式
func FormatJSON(jsonString string) (string, error) {
    if jsonString == "" {
        return "", nil
    }

    var jsonObj interface{}

    // 先尝试解析为通用接口
    err := json.Unmarshal([]byte(jsonString), &jsonObj)
    if err != nil {
        return jsonString, fmt.Errorf("无法解析JSON: %v", err)
    }

    // 格式化为带缩进的JSON
    formattedJSON, err := json.MarshalIndent(jsonObj, "", "    ")
    if err != nil {
        return jsonString, fmt.Errorf("无法格式化JSON: %v", err)
    }

    return string(formattedJSON), nil
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

    var jsonObj interface{}
    if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
        return jsonStr, err
    }

    compacted, err := json.Marshal(jsonObj)
    if err != nil {
        return jsonStr, err
    }

    return string(compacted), nil
}
