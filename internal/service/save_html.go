/**
 * @Time : 2024/6/15 18:19
 * @File : save_html.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description: 保存任意HTML
 */

package service

import (
	"encoding/json"
	"fmt"
	"go-wx-download/config"
	"go-wx-download/pkg/utils"
	"io"
	"net/http"
)

func SaveHtml(r *http.Request, cfg *config.Config, localPath string) *utils.Result {
	if r.Method == "OPTIONS" {
		fmt.Println("OPTIONS")
		return utils.Success("保存成功", "")
	}
	var content = make(map[string]interface{})
	all, err2 := io.ReadAll(r.Body)
	err := json.Unmarshal(all, &content)
	if err != nil || err2 != nil {
		return utils.Failure("保存失败", "")
	}
	head := content["head"].(string)
	body := content["body"].(string)
	utils.WriteAppendFile("text.html", head+body)
	return utils.Success("保存成功", "")
}
