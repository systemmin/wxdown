package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-wx-download/internal/common"
	"go-wx-download/internal/constant"
	"go-wx-download/pkg/utils"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
)

// 解析 home 第一页内容
func parseHome(httpUrl string) (map[string]interface{}, string) {
	// 创建一个 HttpClient 实例，设置超时时间为 60 秒
	client := utils.NewHttpClient(constant.TimeOut)
	// 发送请求
	response, err := client.Request("GET", httpUrl, nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, ""
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response))
	if err != nil {
		log.Println(err)
	}
	collectionName := doc.Find("h2").Text() // 集合名称
	var authorName string                   // 公众号名称
	var list string                         // 文章列表
	var scripts []string                    // JavaScript 名称
	doc.Find("script").Each(func(i int, child *goquery.Selection) {
		scripts = append(scripts, child.Text())
	})
	// 拆分行数组
	rowsArrays := strings.Split(strings.Join(scripts, "\n"), "\n")
	for _, text := range rowsArrays {
		if strings.Contains(text, "html(false)") && strings.Contains(text, "nickname") {
			authorName = strings.Split(text, "\"")[1]
			continue
		}
		if strings.Contains(text, "cgiData.appmsg_list =") {
			list = text[len("cgiData.appmsg_list =") : len(text)-len(".appmsg_list;")]
			break
		}
	}
	dataMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(list), &dataMap)
	if err != nil {
		log.Println("cateStr：", list)
		log.Println("解析响应结果错误：", err)
	}
	title := fmt.Sprintf("%s_%s", authorName, collectionName)
	return dataMap, title
}

func listPageHome(biz, cid, sn, sessionUs, hid string, begin int) map[string]interface{} {
	f := rand.Float64()
	// 创建请求参数对象
	params := make(utils.Params)
	params.Set("r", strconv.FormatFloat(f, 'f', -1, 64)) //    随机数
	params.Set("cid", cid)                               // 分类 id
	params.Set("scene", "18")                            // 场景值
	params.Set("version", "63090a13")                    // 客户端版本
	params.Set("__biz", biz)                             // 公众号 id
	params.Set("begin", strconv.Itoa(begin))             // 开始下标================上一个数组长度
	params.Set("f", "json")                              // 返回数据格式
	params.Set("sn", sn)                                 // 签名或随机字符
	params.Set("session_us", sessionUs)                  // 公众原始 id
	params.Set("action", "appmsg_list")                  // 操作 获取消息列表
	params.Set("hid", hid)                               //
	params.Set("devicetype", "Windows+10+x64")           // 设备类型
	params.Set("ascene", "0")
	params.Set("lang", "zh_CN") // 语言
	params.Set("count", "5")    // 数量
	// 返回格式

	httpUrl := params.ToString(constant.CollectHomeUrl)

	// 创建一个 HttpClient 实例，设置超时时间为 60 秒
	client := utils.NewHttpClient(constant.TimeOut)
	// 发送请求
	response, err := client.Request("POST", httpUrl, nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Println("解析 JSON格式 出错:", err)
	}
	return result
}

// CollectionHome 首页合集
func CollectionHome(urlStr string, path string) ([]common.Collect, error) {
	var collects []common.Collect
	params, err := utils.ParseUrl(urlStr)
	if err != nil {
		return collects, err
	}
	// 文章 id = msgId
	biz := params.Get("__biz")
	// 索引下标
	sessionUs := params.Get("session_us")
	hid := params.Get("hid")
	sn := params.Get("sn")

	// 1、解析首页
	album, title := parseHome(urlStr)
	var category []string
	var categoryList []string
	cateList := album["cate_list"].([]interface{})
	for _, key := range cateList {
		item := key.(map[string]interface{})
		cname := item["cname"].(string)
		category = append(category, cname)
		appmsgList := item["appmsg_list"].([]interface{})
		for _, app := range appmsgList {
			appMsg := app.(map[string]interface{})
			url := appMsg["link"].(string)
			categoryList = append(categoryList, url)
		}
	}
	// 2、开始分页加载数据
	begin := 0   // 开始数组长度
	hasMore := 1 // 0结束 1继续
	for i, name := range category {
		fmt.Printf("分类下标：%d,分类名称：%s\n", i, name)

		// 公众号-集合-分类
		jsName := fmt.Sprintf("%s_%s", title, name)
		// 任务开始文件
		cateFilePath := filepath.Join(path, "task", fmt.Sprintf("%s_%s", jsName, "start.txt"))
		// 任务结束文件
		endFilePath := filepath.Join(path, "task", fmt.Sprintf("%s_%s", jsName, "end.txt"))
		// 采集对象
		collect := common.Collect{Name: jsName, StartPath: cateFilePath, EndPath: endFilePath}
		// 默认分类
		if i == 0 {
			// 页面容量
			begin += len(categoryList)
			collect.List = append(collect.List, categoryList...)
		}
		log.Println("开始下标：", begin)

		for {
			// 分页加载
			result := listPageHome(biz, strconv.Itoa(i), sn, sessionUs, hid, begin)
			if result != nil {
				articleList := result["appmsg_list"].([]interface{}) // 文章列表
				hasMore = int(result["has_more"].(float64))          // 结束标识 0结束 1继续
				for _, article := range articleList {
					articleMap, ok := article.(map[string]interface{})
					if !ok {
						fmt.Println("article is not a map")
					}
					collect.List = append(collect.List, articleMap["link"].(string))
				}
				// 增加分页数
				begin += len(articleList)
			}
			if hasMore == 0 {
				begin = 0 // 重置开始下标
				fmt.Printf("分类：%s 采集完成，路径：%s\n", name, cateFilePath)
				break
			}
		}
		// 追加合集对象到数组
		collects = append(collects, collect)
		utils.WriteAppendFile(cateFilePath, strings.Join(utils.UrlFilter(collect.List), "\n"))
		fmt.Printf("分类：%s 写入完成\n", name)
	}
	return collects, nil
}
