package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-wx-download/internal/constant"
	"regexp"
	"strings"
	"time"
)

type NodeContent struct {
	Node     *goquery.Selection
	Original string // 原始字符串
	Target   string // 修改后字符串
	Type     int    // 0 img 1 bg 2 封面图片 3 音频 4 视频
	Styles   []string
}

// RecursionElement 递归查找图片节点和包含background-image的节点
func RecursionElement(selection *goquery.Selection) []NodeContent {
	var nodes []NodeContent
	// 遍历当前节点下的所有子节点
	selection.Children().Each(func(i int, child *goquery.Selection) {
		// 图片节点处理
		handleImage(child, &nodes)
		// 音频节点处理
		handleAudio(child, &nodes)
		// 处理视频
		handleVideo(child, &nodes)
		// 递归调用查找子节点
		nodes = append(nodes, RecursionElement(child)...)
	})
	return nodes
}

func handleImage(child *goquery.Selection, nodes *[]NodeContent) {
	// 1、检查当前节点是否是图片节点
	if child.Is("img") {
		src, exists := child.Attr("data-src")
		if !exists {
			src, exists = child.Attr("src")
		}
		*nodes = append(*nodes, NodeContent{Node: child, Original: src, Target: "", Type: 0})
	}
	// 2、检查当前节点的 style 属性是否包含 data-lazy-bgimg
	bg, exists := child.Attr("data-lazy-bgimg")
	// css
	var styles []string
	if !exists {
		style, exists := child.Attr("style")
		if exists {
			// 处理样式 替换双引号
			style = strings.ReplaceAll(style, "&quot;", "\"")
			parts := strings.Split(style, ";")
			for _, part := range parts {
				if strings.Contains(part, "background-image") && strings.Contains(part, "url") {
					bg = GetBgImage(part)
					continue
				} else {
					styles = append(styles, part)
				}
			}
		}
	}
	if bg != "" {
		*nodes = append(*nodes, NodeContent{Node: child, Original: bg, Target: "", Type: 1, Styles: styles})
	}

	// 3、检查当前节点的 style 属性是否包含background-image
	var styles2 []string
	style, exists := child.Attr("style")
	if exists && strings.Contains(style, "background-image") && strings.Contains(style, "url") {
		// 处理样式 替换双引号
		style = strings.ReplaceAll(style, "&quot;", "\"")
		parts := strings.Split(style, ";")
		for _, part := range parts {
			if strings.Contains(part, "background-image") {
				bg = GetBgImage(part)
				continue
			} else {
				styles2 = append(styles2, part)
			}
		}
	}
	if bg != "" {
		*nodes = append(*nodes, NodeContent{Node: child, Original: bg, Target: "", Type: 1, Styles: styles2})
	}
	// 4、封面
	//cover, exists := child.Attr("data-cover")
	//if exists {
	//	unescape, err := url.PathUnescape(cover)
	//	if err == nil {
	//		*nodes = append(*nodes, NodeContent{Node: child, Original: unescape, Target: "", Type: 2})
	//	}
	//}
	// 5、svg embed
	if child.Is("embed") {
		src, exists := child.Attr("src")
		if exists {
			*nodes = append(*nodes, NodeContent{Node: child, Original: src, Target: "", Type: 0})
		}

	}
}

func handleAudio(child *goquery.Selection, nodes *[]NodeContent) {
	if child.Is("section") {
		html, err := child.Html()
		if err == nil {
			htmlTrim := strings.Trim(strings.ReplaceAll(html, "\n", ""), " ")
			if strings.HasPrefix(htmlTrim, "<mp-common-mpaudio") && strings.HasSuffix(htmlTrim, "</mp-common-mpaudio>") {
				*nodes = append(*nodes, NodeContent{Node: child, Original: html, Target: "", Type: 3})
			}
			if strings.HasPrefix(htmlTrim, "<mpvoice") && strings.HasSuffix(htmlTrim, "</mpvoice>") {
				*nodes = append(*nodes, NodeContent{Node: child, Original: html, Target: "", Type: 3})
			}
		}
	}
}
func handleVideo(child *goquery.Selection, nodes *[]NodeContent) {
	if child.Is("iframe") {
		attr, exists := child.Attr("data-mpvid")
		if exists {
			*nodes = append(*nodes, NodeContent{Node: child.Parent(), Original: attr, Target: "", Type: 4})
		}
	}
}

// FindMpAudio 递归查找图片节点和包含background-image的节点
func FindMpAudio(selection *goquery.Selection) []NodeContent {
	var nodes []NodeContent
	// 遍历当前节点下的所有子节点
	selection.Children().Each(func(i int, child *goquery.Selection) {
		// 1、检查当前节点是否是图片节点
		if child.Is("section") {
			html, err := child.Html()
			if err == nil {

				htmlTrim := strings.Trim(strings.ReplaceAll(html, "\n", ""), " ")
				if strings.HasPrefix(htmlTrim, "<mp-common-mpaudio") && strings.HasSuffix(htmlTrim, "</mp-common-mpaudio>") {
					nodes = append(nodes, NodeContent{Node: child, Original: html, Target: "", Type: 4})
				}
			}
		}
		// 递归调用查找子节点
		nodes = append(nodes, FindMpAudio(child)...)
	})
	return nodes
}

// ParseScriptVideo 解析 video 数据
func ParseScriptVideo(doc *goquery.Document) ([]string, []string) {
	var scripts []string
	var listData []string
	doc.Find("script").Each(func(i int, child *goquery.Selection) {
		scripts = append(scripts, child.Text())
	})
	// 拆分行数组
	rowsArrays := strings.Split(strings.Join(scripts, "\n"), "\n")
	var videos []string     // 视频地址
	var format []string     // 视频格式 超清 10002 流畅 10004
	var hitVid []string     // 视频 id
	var videoLevel []string // 视频等级
	var coverUrl []string   // 视频封面
	isVideo := false
	for _, v := range rowsArrays {
		if strings.Contains(v, "videoPageInfos") && strings.Contains(v, "[") {
			isVideo = true
			continue
		}
		if strings.Contains(v, "window.__videoPageInfos") {
			isVideo = false
		}
		if isVideo {
			if strings.Contains(v, "video_id") {
				vid := strings.Split(v, "'")
				hitVid = append(hitVid, vid[1])
			}
			if strings.Contains(v, "mp4") {
				videoKey := strings.Split(v, "'")
				videoUrl := strings.ReplaceAll(videoKey[1], "\\x26amp;", "&")
				videoUrl = strings.ReplaceAll(videoUrl, "http://", "https://")
				videos = append(videos, videoUrl)
			}
			if strings.Contains(v, "format_id") {
				formatId := strings.Split(v, "'")
				format = append(format, formatId[1])
			}
			if strings.Contains(v, "video_quality_level") {
				level := strings.Split(v, "'")
				videoLevel = append(videoLevel, level[1])
			}
			if strings.Contains(v, "cover_url") {
				level := strings.Split(v, "'")
				coverUrl = append(coverUrl, level[1])
			}
		}
	}
	if len(videos) == 0 {
		return listData, coverUrl
	}
	group := len(videos) / len(hitVid) // 视频数组长度 / vi数组长度 = 视频可以被分的组数

	for i, key := range videos {
		vIdIndex := i / group // 得到对应的 视频 id
		sprintf := fmt.Sprintf("%s&vid=%s&format_id=%s&support_redirect=0&mmversion=false", key, hitVid[vIdIndex], format[i])
		array, _ := IsValueArray(listData, hitVid[vIdIndex])
		if videoLevel[i] == "3" && len(array) == 0 { // 只保留高清
			listData = append(listData, sprintf)
		}
	}
	return listData, coverUrl
}

// ParseAlbum 解析相册
func ParseAlbum(doc *goquery.Document) []NodeContent {
	var nodes []NodeContent
	compile, _ := regexp.Compile(`https?://[^'\s"]+`)
	doc.Find("script").Each(func(i int, child *goquery.Selection) {
		//picture_page_info_list
		text := child.Text()
		if strings.Contains(text, "picture_page_info_list") {
			findString := compile.FindAllString(text, -1)
			for _, s := range findString {
				if strings.Contains(s, "mmbiz_jpg") {
					nodes = append(nodes, NodeContent{Node: child, Original: s, Target: "", Type: 2})
				}
			}
		}
	})
	return nodes
}

func IsValueArray(array []string, key string) (string, int) {
	for i, item := range array {
		if strings.Contains(item, key) {
			return item, i
		}
	}
	return "", -1
}

// video 对应 iframe；Audio 对应 mp-common-mpaudio

func CreateAudioHTML(title string, src string) string {
	template := fmt.Sprintf("<figure><figcaption class=\"audio_card_title\">%s</figcaption>", title)
	template += fmt.Sprintf("<audio style=\"width: 100%s\" controls src=\"../audios/%s\"></audio></figure>", "%;", src)
	return template
}

func CreateVideoHTML(src string) string {
	template := "<video style='background-color: #000;border-radius: 5px;' width='100%' height='508' controls>"
	template += fmt.Sprintf("<source src=\"../videos/%s\" type='video/mp4'></video>", src)
	return template
}

func ParseScript(doc *goquery.Document) map[string]string {
	var scripts []string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scripts = append(scripts, s.Text())
	})
	join := strings.Join(scripts, "")                 // 拼接
	join = strings.ReplaceAll(join, "__biz", "__BIZ") // 替换干扰内容
	return GetBaseInfo(join)
}

// GetBaseInfo 获取基础信息
func GetBaseInfo(str string) map[string]string {
	// 匹配规则
	keys := constant.Fields
	info := make(map[string]string)

	for _, item := range keys {
		// 构建正则表达式
		regex := regexp.MustCompile(item + `\s*[:=]\s*["']([^"']+)["']`)
		match := regex.FindStringSubmatch(str)
		if len(match) > 1 {
			jsonStr := match[1]
			if item == "biz" {
				// 对biz字段进行base64解码
				buffer, err := base64.StdEncoding.DecodeString(jsonStr)
				if err != nil {
					info[item] = ""
					fmt.Println("Base64 decode error:", err)
				} else {
					info[item+"_base64"] = string(buffer)
					info[item] = jsonStr
				}
			} else {
				info[item] = jsonStr
			}
		} else {
			fmt.Println("没有找到匹配的内容")
		}
	}
	if value, exists := info["createTime"]; exists {
		if len(value) <= 0 {
			info["createTime"] = time.DateTime
		}
	} else {
		info["createTime"] = time.DateTime
	}
	return info
}
