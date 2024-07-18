/**
 * @Time : 2024/6/15 18:16
 * @File : collect.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description: 合集处理控制层，包含标签合集，首页合集
 */

package controller

import (
	"encoding/json"
	"go-wx-download/config"
	"go-wx-download/internal/service"
	"net/http"
)

// Collect 合集
func Collect(w http.ResponseWriter, r *http.Request, cfg *config.Config, path string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	data := service.Collect(r, cfg, path)
	marshal, _ := json.Marshal(data)
	_, err := w.Write(marshal)
	if err != nil {
		return
	}
}
