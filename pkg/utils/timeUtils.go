package utils

import (
	"time"
)

// StrToDate 字符串时间转 time.Time 类型
func StrToDate(strTime string) (time.Time, error) {
	// 使用 Parse 函数将字符串转换为时间类型
	date, err := time.Parse(time.DateOnly, strTime)
	if err != nil {
		//fmt.Println("文件日期转换失败:", err)
		return date, err
	}
	return date, nil
}
