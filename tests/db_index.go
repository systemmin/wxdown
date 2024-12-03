/**
 * @Time : 2024/12/3 11:11
 * @File : db_index.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description:
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// DbToListMap db 文件转 map 数据,
// ddp 文件路径
func DbToListMap(dbp string) map[string]bool {
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

// CheckUrl db 文件转 map 数据,
func CheckUrl(url string, data []map[string]bool) bool {

	for _, datum := range data {
		if datum[url] {
			return true
		}
	}
	return false
}

// LoadDBData db 文件转 map 数据,
func LoadDBData(path string) []map[string]bool {
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	var listMap = make([]map[string]bool, len(dir))
	for _, entry := range dir {
		if entry.Name() == "css" || entry.Name() == "task" {
			continue
		}
		dbpath := filepath.Join(path, entry.Name(), "db", "db.jsonl")
		_, err := os.Stat(dbpath)
		if err == nil {
			toListMap := DbToListMap(dbpath)
			listMap = append(listMap, toListMap)
		}
	}
	return listMap
}

var listData = make([]map[string]bool, 10)

func main() {
	name := "技术最前线"
	html := "2024-11-28-字节起诉前实习生，索赔 800 万.html"
	httpURL := fmt.Sprintf("%s://127.0.0.1:%s/wx/%s/html/%s", "http", "81", url.PathEscape(name), url.PathEscape(html))
	fmt.Println(httpURL)
	// 包含空格的字符串
	str := "Hello World"

	// 使用url.QueryEscape进行URL编码
	encodedStr := url.QueryEscape(str)

	// 打印编码后的字符串
	fmt.Println(encodedStr) // 输出: Hello%20World
}
