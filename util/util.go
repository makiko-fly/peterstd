package util

import (
	"fmt"
	"runtime"
)

// 获取当前代码的位置信息, 默认 skip 为 1
func GetLineNum() string {
	return GetLineNumSkip(1)
}

func GetLineNumSkip(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}
