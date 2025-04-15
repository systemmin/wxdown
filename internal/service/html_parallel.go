package service

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-wx-download/config"
	"go-wx-download/internal/common"
	"go-wx-download/internal/constant"
	"go-wx-download/pkg/down"
	"go-wx-download/pkg/list"
	"go-wx-download/pkg/utils"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 全局变量
var downloaded []string

// HandleDownHTML 下载 html 原始文件
func HandleDownHTML(cfg *config.Config, urlParams *common.UrlParams, host, localPath string, listData []map[string]bool) bool {
	// 开始时间
	start := time.Now()
	var urls []string
	// 处理重复文件
	for _, item := range urlParams.Urls {
		if cfg.Override && len(listData) > 0 && utils.CheckUrl(item, listData) {
			log.Println("已存在，跳过")
			continue
		}
		if cfg.Override && !list.IsEmpty(downloaded) && list.IsExist(downloaded, item) {
			log.Println("已存在，跳过")
			continue
		}
		urls = append(urls, item)
		downloaded = append(downloaded, item)
	}
	// 定义 goroutines 等待组
	var wg sync.WaitGroup
	// 并发数量
	maxConcurrency := utils.Iit(cfg.Thread.Html > 0, cfg.Thread.Html, 10)
	// 创建 sem 通道
	sem := make(chan struct{}, maxConcurrency)
	// 用于存储文件地址的通道
	filePaths := make(chan string, len(urls))
	for i, item := range urls {

		log.Printf("当前下标：[%d]\n", i+1)
		log.Printf("当前资源：[%s]\n", item)
		wg.Add(1)
		// 异步
		go downloadHtml(item, localPath, urlParams.Folder, sem, &wg, filePaths, cfg)
	}
	// 等待所有下载完成
	go func() {
		wg.Wait()
		close(filePaths) // 释放
	}()

	// 开始处理文件转换
	fileFormatConversion(filePaths, localPath, cfg)
	// 结束时间
	end := time.Now()
	duration := end.Sub(start)
	fmt.Println("采集耗时：", duration)
	return true
}

// fileFormatConversion 文件格式转换
func fileFormatConversion(filePaths <-chan string, localPath string, cfg *config.Config) {
	// 创建通信
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup // 等待组
	for file := range filePaths {
		if len(file) <= 0 {
			continue
		}
		sem <- struct{}{} // 占坑
		wg.Add(1)

		go func(f string) {
			defer wg.Done()
			defer func() { <-sem }() // 释放

			// 文件路径解析
			wxName, wxTile, fullTitle := utils.ResolveFilePath(f)
			// pdf
			if cfg.Wkhtmltopdf.Enable {
				protocol := utils.Iif(cfg.Https, "https", "http")
				httpURL := fmt.Sprintf("%s://127.0.0.1:%s/wx/%s/html/%s", protocol, cfg.Port, url.PathEscape(wxName), url.PathEscape(fullTitle))
				output := filepath.Join(localPath, wxName, "pdf", wxTile+".pdf")
				utils.ToPDF(output, httpURL, cfg.Wkhtmltopdf.Path)
			}
			// word
			if cfg.Pandoc.Enable {
				output := filepath.Join(localPath, wxName, "pdf", wxTile+".docx")
				utils.ToWord(f, output, cfg.Pandoc.Path)
			}

		}(file)
	}
	wg.Wait()
	fmt.Println("全部转换结束！")
}

// filePaths 返回文件完整路径
func downloadHtml(urlStr, path, newName string, sem chan struct{}, wg *sync.WaitGroup, filePaths chan string, cfg *config.Config) {
	defer wg.Done()

	// 从信号量中获取一个令牌
	sem <- struct{}{}
	defer func() { <-sem }() // 确保在函数返回时释放信号量令牌

	// 视频自定义 headers
	var headers = make(map[string]string)
	headers["Host"] = "mpvideo.qpic.cn"
	headers["Origin"] = "https://mp.weixin.qq.com"
	headers["Range"] = "bytes=0-"
	headers["Referer"] = "https://mp.weixin.qq.com/"

	// 创建一个 HttpClient 实例，设置超时时间为 60 秒
	client := utils.NewHttpClient(constant.TimeOut)
	// 发送请求
	response, err := client.Request("GET", urlStr, nil, nil)
	if err != nil {
		fmt.Println("创建请求实例时出错:", err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response))
	if err != nil {
		log.Println("NewDocumentFromReader 解析异常:", err)
		return
	}

	// 获取基础信息
	baseInfo := utils.ParseScript(doc)

	// 设置文档可见
	jsContent := doc.Find("#js_content")
	jsContent.SetAttr("style", "visibility: visible")

	// 文章标题
	activityName := doc.Find("#activity-name").Text()
	if len(activityName) == 0 { // 图片轮播集
		activityName = baseInfo["msg_title"]
		fmt.Println(activityName)
	}
	// 图集表示
	isAlbum := false

	// 公众号名称
	jsName := doc.Find("#js_name").Text()
	if len(jsName) == 0 { // 图片轮播集
		jsName = doc.Find(".wx_follow_nickname").Eq(0).Text()
		isAlbum = true
		fmt.Println(jsName)
	}
	jsName = utils.Iif(len(newName) > 0, newName, jsName)

	fmt.Println("公众号名称：", strings.TrimSpace(jsName))
	baseInfo["activity_Name"] = utils.SanitizeFilename(activityName)
	baseInfo["js_name"] = utils.SanitizeFilename(jsName)
	baseInfo["url"] = urlStr

	jsonData, err := json.Marshal(baseInfo)
	if err != nil {
		log.Println("将映射封送至JSON时出错:", err)
		return
	}
	baseInfoStr := strings.ReplaceAll(string(jsonData), "\\u0026", "&")
	fmt.Printf("基础信息：%s\n", baseInfoStr)
	if len(jsName) <= 0 {
		log.Printf("公众号名称为空：%s，未采集到内容\n", jsName)
		return
	}
	// 创建公众号文件夹
	utils.CreateNewFolder(baseInfo["js_name"], path)

	// 文件名命名规则 时间+md5(文件名)+序号+文件类型
	createTime := baseInfo["createTime"][0:10]
	videoUrls, coverUrls := utils.ParseScriptVideo(doc)
	// 获取所有图片遍历并下载
	nodes := utils.RecursionElement(jsContent)
	// 图片合集
	if isAlbum {
		nodes = append(nodes, utils.ParseAlbum(doc)...)
	}

	var wgFile sync.WaitGroup
	// 并发数量
	maxConcurrency := utils.Iit(cfg.Thread.Image > 0, cfg.Thread.Image, 20)
	// 创建 sem 通道
	semFile := make(chan struct{}, maxConcurrency)
	// 创建 base64 通道
	base64String := make(chan string, len(nodes))
	// 创建 set
	sets := make(map[string]string)
	for i, node := range nodes {
		original := node.Original
		types := node.Type
		// 0 1 2 图片 , 3 音频 4 视频
		switch types {
		case 3:
			var voiceEncodeFileID string
			var voiceName string
			audio, err := common.AudioParse(original)
			if err != nil {
				fmt.Println(err)
				voice, _ := common.VoiceParse(original)
				voiceEncodeFileID = voice.VoiceEncodeFileID
				voiceName = voice.Name
			} else {
				voiceEncodeFileID = audio.VoiceEncodeFileID
				voiceName = audio.Name
			}
			// 发布时间->音频名称-->音频序号->音频后缀
			audioName := fmt.Sprintf("%s_%s_%d.%s", createTime, audio.Name, i, "mp3")
			audioPath := filepath.Join(path, baseInfo["js_name"], "audios", audioName)
			wgFile.Add(1)
			down.DownloadFile(constant.AudioPrefix+voiceEncodeFileID, audioPath, nil, semFile, &wgFile)
			// 清空节点内容进行覆盖
			node.Node.SetHtml(utils.CreateAudioHTML(voiceName, audioName))
			break
		case 4:
			vUrl, index := utils.IsValueArray(videoUrls, original)
			if len(vUrl) > 0 {
				parse, err := url.Parse(vUrl)
				if err != nil {
					log.Println(err)
					continue
				}
				name := parse.Path[1:]
				// 发布时间->序号->名称
				videoName := fmt.Sprintf("%s_%d_%s", createTime, index, name)
				videoPath := filepath.Join(path, baseInfo["js_name"], "videos", videoName)
				wgFile.Add(1)
				down.DownloadFile(vUrl, videoPath, headers, semFile, &wgFile)
				// 清空节点内容进行覆盖
				node.Node.SetHtml(utils.CreateVideoHTML(videoName))
				// 下载封面
				cUrl := coverUrls[index]
				all := strings.ReplaceAll(videoName, "mp4", "jpeg")
				imgPath := filepath.Join(path, baseInfo["js_name"], "images", all)
				wgFile.Add(1)
				down.DownloadFile(cUrl, imgPath, nil, semFile, &wgFile)
			}
			break
		default:
			// 文件名 MD5 哈希值
			hash := md5.New()
			hash.Write([]byte(original))
			md5Sum := hash.Sum(nil)

			// 处理图片后缀
			suffix := utils.GetUrlParams(original, "wx_fmt")
			if len(suffix) == 0 {
				suffix = utils.GetSuffix(original)
			}
			fileName := fmt.Sprintf("%s_%x_%d.%s", createTime, md5Sum, i, suffix)
			imageJoin := filepath.Join(path, baseInfo["js_name"], "images", fileName)
			resetSrc := "../images/" + fileName
			// 图片转Base
			if cfg.Base64 {
				// 添加线程计数
				wgFile.Add(1)
				resetSrc = down.ImageToBase64(original, suffix, nil, semFile, &wgFile)
			} else {
				// 未下载调用线程进行下载
				if len(sets[string(md5Sum)]) == 0 {
					// 添加线程计数
					wgFile.Add(1)
					down.DownloadFile(original, imageJoin, nil, semFile, &wgFile)
					sets[string(md5Sum)] = fileName
				} else {
					// 重新赋值已下载内容
					resetSrc = "../images/" + sets[string(md5Sum)]
				}
				before, b := strings.CutSuffix(resetSrc, "webp") // 转 PDF 有问题图片格式处理
				if b {
					resetSrc = fmt.Sprintf("%sjpeg", before)
				}
			}
			// 重置属性值
			if node.Type == 0 {
				node.Node.SetAttr("src", resetSrc)
				node.Node.SetAttr("data-src", resetSrc)
			} else if node.Type == 1 {
				node.Node.SetAttr("data-lazy-bgimg", resetSrc)
				styles := node.Styles
				// 将 css 背景图片添加到数组中
				styles = append(styles, "background-image: url("+resetSrc+")")
				// 拼接
				join := strings.Join(styles, ";")
				// 重新设置 style
				node.Node.SetAttr("style", join)
			} else if node.Type == 2 {
				text := node.Node.Text()
				replace := strings.Replace(text, node.Target, resetSrc, -1)
				node.Node.SetHtml(replace)
			}
		}

	}
	// 兼容旧版本视频
	video := utils.ParseScriptOldVideo(doc)
	if len(video) > 0 {
		// 发布时间->序号->名称
		videoName := fmt.Sprintf("%s_%d_%s.mp4", createTime, 0, baseInfo["activity_Name"])
		videoPath := filepath.Join(path, baseInfo["js_name"], "videos", videoName)
		wgFile.Add(1)
		down.DownloadFile(video, videoPath, headers, semFile, &wgFile)
	}
	// 等待所有下载完成
	go func() {
		wgFile.Wait()
		close(base64String)
	}()

	// 定义 css 路径
	var css []string
	if isAlbum {
		css = append(css, "index.m0jn22vy4f03b36c.css", "qqmail_tpl_vite_entry.m0jn22vyffac437b.css")
	} else {
		css = append(css,
			"index.lyptmz0d196f5b68.css",
			"cover_next.lyptmz0d8abb784e.css",
			"interaction.lyptmz0d9570c58b.css",
			"qqmail_tpl_vite_entry.lyptmz0da92f2c62.css",
			"tencent_portfolio_light.lyptmz0d0cd74df8.css",
			"weui.min.css",
		)
	}

	// 修改 link 引入文件路径
	doc.Find("link[rel='stylesheet']").Each(func(i int, el *goquery.Selection) {
		href, exists := el.Attr("href")
		if exists && len(href) > 0 {
			// 重置 href 路径
			if i < len(css) {
				el.SetAttr("href", "../../css/"+css[i])
			}
		}
	})
	// 公众号二维码加载
	doc.Find("#js_pc_qr_code").Remove()
	// html 写入本地
	html, err := doc.Html()
	if err != nil {
		fmt.Println("生成 HTML 内容时出错:", err)
		return
	}
	// 采集页面信息 URL 替换
	//protocol := utils.Iif(cfg.Https, "https", "http")
	html = strings.ReplaceAll(html, "https://badjs.weixinbridge.com", "")
	// 字体格式处理，防止中文乱码 @see https://github.com/systemmin/wxdown/issues/16#issuecomment-2726465545
	html = strings.ReplaceAll(html, "仿宋", "FangSong")

	fileName := fmt.Sprintf("%s-%s.html", baseInfo["createTime"][0:10], baseInfo["activity_Name"])
	join := filepath.Join(path, baseInfo["js_name"], "html", fileName)
	htmlErr := os.WriteFile(join, []byte(html), 0644)
	if htmlErr != nil {
		fmt.Println("写入 HTML 文件时出错:", htmlErr)
		return
	}

	// 打开文件，如果文件不存在则创建，使用 os.O_APPEND 标志表示追加内容
	file, err := os.OpenFile(filepath.Join(path, baseInfo["js_name"], "db", "db.jsonl"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 写入内容
	if _, err := file.WriteString(baseInfoStr + "\n"); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	// 发送文件路径到通道
	filePaths <- join

	fmt.Println("写入成功:", join)
}
