/**
 * @Time : 2024/12/3 11:11
 * @File : db_index.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description: 加载 db 文件，建立检索索引
 */

package utils

import (
	"encoding/json"
	"fmt"
	"go-wx-download/internal/constant"
	"go-wx-download/pkg/list"
	"os"
	"path/filepath"
	"strings"
)

// DbToListMap db 文件转 map 数据,
// ddp 文件路径
func dbToListMap(dbp string) map[string]bool {
	var dataMap = make(map[string]bool)
	// 1、加载 db 文件
	file, err := os.ReadFile(dbp)
	if err != nil {
		return dataMap
	}
	// 2、拆分+拼接 json
	lines := strings.Split(string(file), "\n")
	joinMap := strings.Join(lines, ",")
	jsonContent := fmt.Sprintf("[%s]", joinMap[:len(joinMap)-1])
	// 3、遍历去重
	var listData []map[string]any
	err = json.Unmarshal([]byte(jsonContent), &listData)
	if err != nil {
		return dataMap
	}
	for _, data := range listData {
		key := data["url"].(string)
		dataMap[key] = true
	}
	return dataMap
}

// CheckUrl 检查 URL 是否存在
func CheckUrl(url string, data []map[string]bool) bool {
	for _, datum := range data {
		if len(datum) == 0 {
			continue
		}
		if datum[url] {
			return true
		}
	}
	return false
}

// LoadDBData 加载数据
func LoadDBData(path string) []map[string]bool {
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	var listMap = make([]map[string]bool, len(dir))
	for _, entry := range dir {
		if list.IsExist(constant.ExcludeFolder, entry.Name()) {
			continue
		}
		dbpath := filepath.Join(path, entry.Name(), "db", "db.jsonl")
		_, err := os.Stat(dbpath)
		if err == nil {
			toListMap := dbToListMap(dbpath)
			fmt.Println(dbpath, ":", len(toListMap))
			listMap = append(listMap, toListMap)
		}
	}
	return listMap
}
