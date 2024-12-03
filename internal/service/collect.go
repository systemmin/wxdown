package service

import (
	"fmt"
	"go-wx-download/config"
	"go-wx-download/internal/common"
	"go-wx-download/pkg/list"
	"go-wx-download/pkg/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Collect(r *http.Request, cfg *config.Config, localPath string, listData []map[string]bool) *utils.Result {
	host := r.Host
	params, err := common.GetParams(r)
	if err != nil {
		return utils.Failure("参数解析异常", err)
	}
	fmt.Println("Folder：", params.Folder)
	fmt.Println("URL：", params.Urls)

	if list.IsEmpty(params.Urls) {
		return utils.Failure("地址为空", "")
	}
	url := params.Urls[0]

	if err = utils.IsURL(url); err != nil {
		return utils.Failure("地址不符和规范", err)
	}

	if !strings.Contains(url, "appmsgalbum") && !strings.Contains(url, "homepage") {
		return utils.Failure("不是合集地址", url)
	}

	var collects []common.Collect
	if strings.Contains(url, "appmsgalbum") {
		album, err := CollectionAlbum(url, localPath)
		if err != nil {
			return utils.Failure("合集解析异常", err)
		}
		collects = append(collects, album)
	} else if strings.Contains(url, "homepage") {
		home, err := CollectionHome(url, localPath)
		if err != nil {
			return utils.Failure("合集解析异常", err)
		}
		collects = home
	}
	// 开始时间
	start := time.Now()
	log.Println("任务开始时间：", start.Format(time.DateTime))
	for _, collect := range collects {
		params.Urls = collect.List   // 覆盖 URL
		if len(params.Folder) == 0 { // 覆盖 名称
			params.Folder = collect.Name
		}
		html := HandleDownHTML(cfg, &params, host, localPath, listData)
		if html {
			// 拷贝任务文件
			utils.CopyFile(collect.EndPath, collect.StartPath)
			// 删除任务文件
			os.Remove(collect.StartPath)
		}
	}
	end := time.Now()
	duration := end.Sub(start)
	log.Println("任务结束时间：", end.Format(time.DateTime))
	fmt.Println("任务采集耗时：", duration)
	return utils.Success("操作成功", "")
}
