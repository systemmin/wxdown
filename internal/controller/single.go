/**
 * @Time : 2024/6/15 18:16
 * @File : single.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description: 单个链接，多个链接，提交采集，支持 get post put 等多种请求方式
 */

package controller

import (
	"encoding/json"
	"go-wx-download/config"
	"go-wx-download/internal/service"
	"net/http"
)

// Gather w 响应对象 r 请求对象 path 本地文件根路径 port 端口
func Gather(w http.ResponseWriter, r *http.Request, cfg *config.Config, path string, listData []map[string]bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	data := service.Gather(r, cfg, path, listData)
	marshal, _ := json.Marshal(data)
	_, err := w.Write(marshal)
	if err != nil {
		return
	}
}
