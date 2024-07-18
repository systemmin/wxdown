/**
 * @Time : 2024/6/15 18:16
 * @File : collect.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description: web 接口统一返回结构体
 */

package utils

import (
	"fmt"
	"time"
)

// Result 定义 Result 结构体
type Result struct {
	Code      int         `json:"code"`
	Success   bool        `json:"success"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功
func Success(msg string, data interface{}) *Result {
	return &Result{
		Code:      200,
		Success:   true,
		Msg:       msg,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

// Failure 失败
func Failure(msg string, data interface{}) *Result {
	return &Result{
		Code:      500,
		Success:   false,
		Msg:       msg,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	}
}

func NotLogin(msg string) *Result {
	return &Result{
		Code:      401,
		Success:   false,
		Msg:       msg,
		Data:      nil,
		Timestamp: time.Now().UnixMilli(),
	}
}

// 实现 toString 方法，用于打印 Result 的内容
func (r *Result) String() string {
	return fmt.Sprintf("Result{code=%d, success=%t, msg='%s', data=%v}", r.Code, r.Success, r.Msg, r.Data)
}
