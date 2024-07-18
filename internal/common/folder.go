/**
 * @Time : 2024/6/15 18:22
 * @File : folder.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description: 文件目录结构体
 */

package common

import "time"

type Folder struct {
	Name     string    `json:"name"`     // 名称
	ModTime  time.Time `json:"modTime"`  // 修改时间
	Link     string    `json:"link"`     // 访问链接
	Path     string    `json:"path"`     // 本地路径
	CteTime  string    `json:"cteTime"`  // 创建时间
	Original string    `json:"original"` // 原始名称
}
