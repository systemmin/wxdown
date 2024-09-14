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
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func listPageAlbum(biz, albumId, mid, idx string) map[string]interface{} {
	// 创建请求参数对象
	params := make(utils.Params)
	params.Set("__biz", biz)                   // 公众号标识
	params.Set("action", "getalbum")           // 操作：获取专辑合集
	params.Set("album_id", albumId)            // 合集 id
	params.Set("count", "10")                  // 数量
	params.Set("begin_msgid", mid)             // 消息 id ； 下一个列表取上一页列表的 msgId
	params.Set("begin_itemidx", idx)           // 消息索引
	params.Set("devicetype", "Windows;10;x64") // 设备类型
	params.Set("clientversion", "63090a13")    // 客户端版本
	params.Set("x5", "0")                      // 是否为微信浏览器
	params.Set("f", "json")                    // 返回格式

	httpUrl := params.ToString(constant.CollectAlbumUrl)
	log.Println(httpUrl)

	client := utils.NewHttpClient(constant.TimeOut)
	response, err := client.Request("GET", httpUrl, nil, nil)
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

// 解析合集页面dom
func parseAlbum(httpUrl string) ([]string, string) {
	log.Println("分页地址：", httpUrl)

	client := utils.NewHttpClient(constant.TimeOut)
	response, err := client.Request("GET", httpUrl, nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, ""
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response))
	if err != nil {
		log.Println(err)
	}
	// 视频合集标识
	flagVideo := false

	text := doc.Find("#js_tag_name").Text()
	authorName := doc.Find(".album__author-name").First().Text()
	if len(authorName) == 0 {
		authorName = doc.Find(".video-album_account-name").First().Text()
		if len(authorName) > 0 {
			flagVideo = true
		}
	}
	title := fmt.Sprintf("%s_%s", authorName, text)
	var urls []string
	log.Println("解析 dom 文章地址")
	doc.Find(".js_album_list li").Each(func(i int, selection *goquery.Selection) {
		attr, exists := selection.Attr("data-link")
		if exists {
			urls = append(urls, attr)
			log.Println(i, ":", attr)
		}
	})
	// 正则匹配，解析视频合集
	if flagVideo {
		log.Println("解析视频合集地址")
		compile, _ := regexp.Compile(`https?://[^'\s"]+`)
		doc.Find("script").Each(func(i int, child *goquery.Selection) {
			script := child.Text()
			if strings.Contains(script, "videoList") {
				findString := compile.FindAllString(script, -1)
				for _, s := range findString {
					if strings.Contains(s, "__biz") {
						attr := strings.ReplaceAll(s, "&amp;amp;", "&")
						urls = append(urls, attr)
						log.Println(i, ":", attr)
					}
				}
			}
		})
	}
	return urls, title
}

// CollectionAlbum 标签合集
func CollectionAlbum(urlStr string, path string) (common.Collect, error) {
	params, err := utils.ParseUrl(urlStr)
	if err != nil {
		log.Println(err)
	}
	//文章 id = msgId
	biz := params.Get("__biz")
	// 索引下标
	albumId := params.Get("album_id")

	// 1、解析首页
	album, title := parseAlbum(urlStr)
	title = utils.SanitizeFilename(title)
	// 边读边写
	join := filepath.Join(path, "task", fmt.Sprintf("%s_%s", title, "start.txt"))
	endPath := filepath.Join(path, "task", fmt.Sprintf("%s_%s", title, "end.txt"))
	log.Println("写入任务地址:", join)
	utils.WriteAppendFile(join, strings.Join(album, "\n"))
	// 2、获取下页需要的参数
	last := album[len(album)-1]
	params, err = utils.ParseUrl(last)
	if err != nil {
		log.Println(err)
	}
	// 文章 id = msgId
	mid := params.Get("mid")
	// 索引下标
	idx := params.Get("idx")
	log.Println(biz, albumId, mid, idx)
	continueFlag := "1"
	for {
		// 分页加载
		result := listPageAlbum(biz, albumId, mid, idx)
		marshal, _ := json.Marshal(result)
		//  {"base_resp":{"exportkey_token":"","ret":0},"getalbum_resp":{"base_info":{"is_first_screen":"0"},"continue_flag":"0","reverse_continue_flag":"1"}}
		log.Println(string(marshal))
		if result != nil {
			resp := result["getalbum_resp"].(map[string]interface{})
			// 结果最后一条直接返回对象
			articleList, ok := resp["article_list"].([]interface{})
			// continue_flag 0 结束1 继续,2024年9月14日18:14:50 处理 结束表示为0时，可能依然存在结果
			continueFlag = resp["continue_flag"].(string)
			if continueFlag == "0" && !ok {
				break
			}
			var urls []string
			if ok {
				for _, article := range articleList {
					articleMap := article.(map[string]interface{})
					// 文章 id = msgId
					mid = articleMap["msgid"].(string)
					// 索引下标
					idx = articleMap["itemidx"].(string)
					fmt.Println(mid, idx)
					urls = append(urls, articleMap["url"].(string))
				}
			} else {
				articleMap := resp["article_list"].(map[string]interface{})
				// 文章 id = msgId
				mid = articleMap["msgid"].(string)
				// 索引下标
				idx = articleMap["itemidx"].(string)
				fmt.Println(mid, idx)
				urls = append(urls, articleMap["url"].(string))
			}
			utils.WriteAppendFile(join, "\n"+strings.Join(urls, "\n"))
		}

	}
	log.Println("合集地址采集完成")
	file, err := os.ReadFile(join)
	if err != nil {
		return common.Collect{}, err
	}
	list := strings.Split(string(file), "\n")
	return common.Collect{List: list, StartPath: join, EndPath: endPath, Name: title}, nil
}
