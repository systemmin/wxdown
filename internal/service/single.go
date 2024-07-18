package service

import (
	"fmt"
	"go-wx-download/config"
	"go-wx-download/internal/common"
	"go-wx-download/pkg/utils"
	"log"
	"net/http"
)

// Gather 单个链接采集，多个链接采集
func Gather(r *http.Request, cfg *config.Config, localPath string) *utils.Result {
	host := r.Host
	params, err := common.GetParams(r)
	if err != nil {
		return utils.Failure("参数解析异常", err)
	}
	fmt.Println("自定义文件名：", params.Folder)
	fmt.Println("采集URL：", params.Urls)
	// =================== URL 过滤 去重
	results := utils.RemoveDuplicates(utils.UrlFilter(params.Urls))
	if len(results) == 0 {
		return utils.Failure("地址不符和规范", "")
	}
	log.Printf("文章数量：%d\n", len(results))
	HandleDownHTML(cfg, &params, host, localPath)
	return utils.Success("操作成功", "")
}
