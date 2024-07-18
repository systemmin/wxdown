/**
 * @Time : 2024/6/15 18:19
 * @File : folder.go
 * @Software: wxdown
 * @Author : Mr.Fang
 * @Description: 本地文件处理业务层
 */

package service

import (
	"fmt"
	"go-wx-download/internal/common"
	"go-wx-download/internal/constant"
	"go-wx-download/pkg/list"
	"go-wx-download/pkg/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Folder 目录下的文件 。folder 目录、fType 文件类型、path 路径
func Folder(localPath string, r *http.Request) *utils.Result {
	folder := r.URL.Path[len("/ats/"):]
	folderType := ""
	if len(folder) > 0 {
		params := strings.Split(folder, "/")
		folder = params[0]
		if len(params) >= 2 {
			folderType = params[1]
		}
	}
	var folders []common.Folder
	if len(folder) == 0 { // 查目录
		folders = getFolders(localPath)

	} else if len(folderType) > 0 { // 查文件
		folders = getFolderDetail(localPath, folder, folderType)
	}
	return utils.Success("查询结果", folders)
}

// Open 打开指定目录
func Open(localPath string, r *http.Request) *utils.Result {
	path := r.URL.Path
	pathSegments := strings.Split(path, "/")
	lastSegment := pathSegments[len(pathSegments)-1]
	// 环境 退出 命令
	var cmd string
	switch runtime.GOOS {
	case "windows":
		cmd = "start"
	case "darwin":
		cmd = "open"
	default:
		cmd = "open"
	}
	arg := filepath.Join(localPath, lastSegment)
	utils.ExecuteCmd(fmt.Sprintf("%s %s", cmd, arg))
	return utils.Success("操作成功", "")
}

// getFolders 获取公众号所有目录列表
func getFolders(path string) []common.Folder {
	var folders []common.Folder
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Printf("%s 路径不存在\n", path)
		return folders
	}
	for _, file := range dir {
		if file.IsDir() {
			name := file.Name()
			info, _ := file.Info()
			if list.IsExist(constant.ExcludeFolder, name) {
				continue
			}
			link, _ := url.JoinPath("/", "ast", name)
			folders = append(folders, common.Folder{Name: name, ModTime: info.ModTime(), Path: filepath.Join(path, name), Link: link})
		}
	}
	// 对文件列表按照修改时间进行排序，最近创建排在最前面
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].ModTime.After(folders[j].ModTime)
	})
	return folders
}

// getFolderDetail 获取公众号所有目录列表
func getFolderDetail(path string, folder string, fType string) []common.Folder {
	var folders []common.Folder
	join := filepath.Join(path, folder, fType)
	dir, err := os.ReadDir(join)
	if err != nil {
		log.Printf("%s 路径不存在\n", join)
		return folders
	}
	for _, file := range dir {
		if !file.IsDir() {
			name := file.Name()
			info, _ := file.Info()
			if name == ".DS_Store" || name == "css" {
				continue
			}
			timeStr := name[0:10]                        // 截取文章发布时间 := name[0:10]                // 截取文章发布时间
			publishTime, err := utils.StrToDate(timeStr) // 格式 time 对象
			if err != nil {
				publishTime = info.ModTime()
				timeStr = time.Time.Format(info.ModTime(), "2006-01-02 15:04:05")
			}
			var original string
			if fType == "html" {
				original = name[11 : len(name)-len(".html")]
			} else {
				original = name[11 : len(name)-4]
			}
			link, _ := url.JoinPath("/wx", url.PathEscape(folder), fType, url.PathEscape(name))
			folders = append(folders, common.Folder{Name: name, ModTime: publishTime, CteTime: timeStr, Original: original, Path: filepath.Join(join, name), Link: link})
		}
	}
	// 对文件列表按照修改时间进行排序，最近创建排在最前面
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].ModTime.After(folders[j].ModTime)
	})
	return folders
}
