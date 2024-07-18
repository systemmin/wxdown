/**
 * @Time : 2024/6/15 18:16
 * @File : folder.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description:
 */

package controller

import (
	"encoding/json"
	"go-wx-download/internal/service"
	"net/http"
)

// Folder 本地文件目录
func Folder(w http.ResponseWriter, r *http.Request, localPath string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	data := service.Folder(localPath, r)
	marshal, _ := json.Marshal(data)
	_, err := w.Write(marshal)
	if err != nil {
		return
	}
}

// Open 打开本地文件夹
func Open(w http.ResponseWriter, r *http.Request, localPath string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	data := service.Open(localPath, r)
	marshal, _ := json.Marshal(data)
	_, err := w.Write(marshal)
	if err != nil {
		return
	}
}
