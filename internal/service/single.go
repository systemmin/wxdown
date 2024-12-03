package service

import (
	"fmt"
	"go-wx-download/config"
	"go-wx-download/internal/common"
	"go-wx-download/pkg/utils"
	"net/http"
)

// Gather 单个链接采集，多个链接采集
func Gather(r *http.Request, cfg *config.Config, localPath string, listData []map[string]bool) *utils.Result {
	host := r.Host
	params, err := common.GetParams(r)
	if err != nil {
		return utils.Failure("参数解析异常", err)
	}
	fmt.Println("Folder：", params.Folder)
	fmt.Println("URL：", params.Urls)
	fmt.Println("数量：", len(params.Urls))
	// =================== URL 过滤 去重
	results := utils.RemoveDuplicates(utils.UrlFilter(params.Urls))
	if len(results) == 0 {
		return utils.Failure("地址不符和规范", "")
	}
	// 覆盖
	params.Urls = results
	HandleDownHTML(cfg, &params, host, localPath, listData)
	return utils.Success("操作成功", "")
}
