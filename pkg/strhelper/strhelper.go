package strhelper

import (
	"strconv"
	"strings"
)

func FloatSliceToString(slice []float64) string {
	// 用于存储转换后的字符串
	strs := make([]string, len(slice))
	for i, v := range slice {
		// 将浮点数转换为字符串，并保留一位小数
		strs[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	// 使用","连接所有字符串，并添加方括号
	return "[" + strings.Join(strs, ",") + "]"
}
