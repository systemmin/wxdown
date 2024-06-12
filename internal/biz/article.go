package biz

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-wx-download/internal/common"
	"go-wx-download/internal/constant"
	"go-wx-download/pkg/utils"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Folder 文件夹对象
type Folder struct {
	Name     string    `json:"name"`     // 名称
	ModTime  time.Time `json:"modTime"`  // 修改时间
	Link     string    `json:"link"`     // 访问链接
	Path     string    `json:"path"`     // 本地路径
	CteTime  string    `json:"cteTime"`  // 创建时间
	Original string    `json:"original"` // 原始名称
}

type Collect struct {
	List      []string // 合计所有 URL 地址
	StartPath string   // 未处理任务本地存储路径
	EndPath   string   // 已处理任务存储路径
	Name      string   // 合计名称
}

// GetFolders 获取公众号所有目录列表
func GetFolders(path string) []Folder {
	var folders []Folder
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Printf("%s 路径不存在\n", path)
		return folders
	}
	for _, file := range dir {
		if file.IsDir() {
			name := file.Name()
			info, _ := file.Info()
			if name == ".DS_Store" || name == "css" || name == "task" {
				continue
			}
			link, _ := url.JoinPath("/", "ast", name)
			folders = append(folders, Folder{Name: name, ModTime: info.ModTime(), Path: filepath.Join(path, name), Link: link})
		}
	}
	// 对文件列表按照修改时间进行排序，最近创建排在最前面
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].ModTime.After(folders[j].ModTime)
	})
	return folders
}

// GetFolderDetail 获取公众号所有目录列表
func GetFolderDetail(path string, folder string, fType string) []Folder {
	var folders []Folder
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
			folders = append(folders, Folder{Name: name, ModTime: publishTime, CteTime: timeStr, Original: original, Path: filepath.Join(join, name), Link: link})
		}
	}
	// 对文件列表按照修改时间进行排序，最近创建排在最前面
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].ModTime.After(folders[j].ModTime)
	})
	return folders
}

// ParseScript 解析 script 标签内容
func parseScript(doc *goquery.Document) map[string]string {
	var scripts []string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scripts = append(scripts, s.Text())
	})
	join := strings.Join(scripts, "")                 // 拼接
	join = strings.ReplaceAll(join, "__biz", "__BIZ") // 替换干扰内容
	return utils.GetBaseInfo(join)
}

// DownloadHtml 下载公众号原始 html 内容
func DownloadHtml(urlStr string, path string, newName string, sem chan struct{}, wg *sync.WaitGroup, filePaths chan string, thread int) {
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

	// 创建请求客户端对象
	client := &http.Client{
		Timeout: 60 * time.Second, // 设置超时时间为30秒
	}

	// 创建一个请求实例
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		fmt.Println("创建请求实例时出错:", err)
		return
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	// 创建一个HTTP客户端并发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时出错:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("NewDocumentFromReader 解析异常:", err)
		return
	}

	// 获取基础信息
	baseInfo := parseScript(doc)

	// 设置文档可见
	jsContent := doc.Find("#js_content")
	jsContent.SetAttr("style", "visibility: visible")

	// 文章标题
	activityName := doc.Find("#activity-name").Text()
	// 公众号名称
	jsName := doc.Find(".profile_nickname").Eq(0).Text()
	jsName = utils.Iif(len(newName) > 0, newName, jsName)

	fmt.Println("公众号名称：" + jsName)
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

	var wgFile sync.WaitGroup
	// 并发数量
	maxConcurrency := utils.Iit(thread > 0, thread, 10)
	// 创建 sem 通道
	semFile := make(chan struct{}, maxConcurrency)
	// 创建 set
	sets := make(map[string]string)
	for i, node := range nodes {
		original := node.Original
		types := node.Type
		target := node.Target
		// 0 1 2 图片 , 3 音频 4 视频
		switch types {
		case 3:
			var voiceEncodeFileID string
			var voiceName string
			if target == "a" {
				audio, _ := common.AudioParse(original)
				voiceEncodeFileID = audio.VoiceEncodeFileID
				voiceName = audio.Name
			}
			if target == "v" {
				audio, _ := common.VoiceParse(original)
				voiceEncodeFileID = audio.VoiceEncodeFileID
				voiceName = audio.Name
			}
			// 发布时间->音频名称-->音频序号->音频后缀
			audioName := fmt.Sprintf("%s_%s_%d.%s", createTime, voiceName, i, "mp3")
			audioPath := filepath.Join(path, baseInfo["js_name"], "audios", audioName)
			wgFile.Add(1)
			go utils.DownloadFile(constant.AudioPrefix+voiceEncodeFileID, audioPath, nil, semFile, &wgFile)
			// 清空节点内容进行覆盖
			text := node.Node.Text()
			node.Node.SetHtml(utils.CreateAudioHTML(voiceName, audioName, text))
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
				go utils.DownloadFile(vUrl, videoPath, headers, semFile, &wgFile)
				// 清空节点内容进行覆盖
				node.Node.SetHtml(utils.CreateVideoHTML(videoName))
				// 下载封面
				cUrl := coverUrls[index]
				all := strings.ReplaceAll(videoName, "mp4", "jpeg")
				imgPath := filepath.Join(path, baseInfo["js_name"], "images", all)
				wgFile.Add(1)
				go utils.DownloadFile(cUrl, imgPath, nil, semFile, &wgFile)
			}
			break
		default:
			// 计算图片内容的 MD5 哈希值
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
			// 未下载调用线程进行下载
			if len(sets[string(md5Sum)]) == 0 {
				// 添加线程计数
				wgFile.Add(1)
				go utils.DownloadFile(original, imageJoin, nil, semFile, &wgFile)
				sets[string(md5Sum)] = fileName
			} else {
				// 重新赋值已下载内容
				resetSrc = "../images/" + sets[string(md5Sum)]
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
				styles := node.Styles
				// 将 css 背景图片添加到数组中
				last := styles[len(styles)-1]
				styles = append(styles[:len(styles)-1], "background: url("+resetSrc+") "+last)
				// 拼接
				join := strings.Join(styles, ";")
				// 重新设置 style
				node.Node.SetAttr("style", join)
			}
		}

	}
	// 等待所有下载完成
	go func() {
		wgFile.Wait()
	}()

	// 定义 css 路径
	css := [...]string{
		"index.lw5sif657992cef9.css",
		"wxwork_hidden.lw5sif65649ec4c3.css",
		"like_and_share.lw5sif65e3b0c442.css",
		"interaction.lw5sif653ab748dd.css",
		"article_bottom_bar.lw5sif6588987ef0.css",
		"controller.lw5sif657f351e8c.css",
		"qqmail_tpl_vite_entry.lw5sif65375661d4.css",
		"tencent_portfolio_light.lw5sif653b934b9b.css",
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
	html = strings.ReplaceAll(html, "https://badjs.weixinbridge.com", "http://localhost:81")

	fileName := fmt.Sprintf("%s-%s.html", baseInfo["createTime"][0:10], baseInfo["activity_Name"])
	join := filepath.Join(path, baseInfo["js_name"], "html", fileName)
	htmlErr := os.WriteFile(join, []byte(html), 0644)
	if htmlErr != nil {
		fmt.Println("写入 HTML 文件时出错:", htmlErr)
		return
	}

	// 打开文件，如果文件不存在则创建，使用os.O_APPEND标志表示追加内容
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

// =============================================标签合计
// 获取专辑列表
// biz 公众号 albumId 专辑 mid 消息 idx 索引
func listPageAlbum(biz, albumId, mid, idx string) map[string]interface{} {
	baseUrl := "https://mp.weixin.qq.com/mp/appmsgalbum"
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

	httpUrl := params.ToString(baseUrl)

	log.Println(httpUrl)
	// 创建请求客户端对象
	client := &http.Client{
		Timeout: 60 * time.Second, // 设置超时时间为30秒
	}

	// 创建一个请求实例
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		fmt.Println("创建请求实例时出错:", err)
		return nil
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	for s, strings := range req.Header {
		fmt.Println(s, strings)
	}
	log.Println(req.URL)

	// 创建一个HTTP客户端并发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时出错:", err)
		return nil
	}
	defer resp.Body.Close()

	result := make(map[string]interface{})
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体时出错:", err)
		return nil
	}
	fmt.Println(string(all))
	err = json.Unmarshal(all, &result)
	if err != nil {
		fmt.Println("解析 JSON格式 出错:", err)
	}
	return result
}

// 解析合集页面dom
func parseAlbum(httpUrl string) ([]string, string) {

	log.Println("分页地址：", httpUrl)
	// 创建请求客户端对象
	client := &http.Client{
		Timeout: 60 * time.Second, // 设置超时时间为30秒
	}

	// 创建一个请求实例
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		fmt.Println("创建请求实例时出错:", err)
		return nil, ""
	}
	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	// 创建一个HTTP客户端并发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时出错:", err)
		return nil, ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}
	text := doc.Find("#js_tag_name").Text()
	authorName := doc.Find(".album__author-name").First().Text()
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
	return urls, title
}

// StartCollectionAlbum urlStr 合集地址
func StartCollectionAlbum(urlStr string, path string) (Collect, error) {
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
	// 边读编写
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
		//log.Println("睡一会儿 zzzz，莫慌 5 秒中………………")
		//time.Sleep(time.Second * 5)
		result := listPageAlbum(biz, albumId, mid, idx)
		if result != nil {
			resp := result["getalbum_resp"].(map[string]interface{})
			// continue_flag 0 结束1 继续
			continueFlag = resp["continue_flag"].(string)
			// 结果最后一条直接返回对象
			articleList, ok := resp["article_list"].([]interface{})
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
		//continue_flag 0 结束1 继续
		if continueFlag == "0" {
			break
		}
	}
	log.Println("合集地址采集完成")
	file, err := os.ReadFile(join)
	if err != nil {
		return Collect{}, err
	}

	split := strings.Split(string(file), "\n")
	return Collect{List: split, StartPath: join, EndPath: endPath, Name: title}, nil
}

// =============================================首页合计

// cid 分类序号 sn 签名字符串 session_us 公众号原始 id idx 索引
// has_more 0 结束
func listPageHome(biz, cid, sn, sessionUs, hid string, begin int) map[string]interface{} {
	rand.NewSource(1)
	f := rand.Float64()
	baseUrl := "https://mp.weixin.qq.com/mp/homepage"
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

	httpUrl := params.ToString(baseUrl)

	// 创建一个 HttpClient 实例，设置超时时间为 60 秒
	client := utils.NewHttpClient(60 * time.Second)
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

// 解析 home 第一页内容
func parseHome(httpUrl string) (map[string]interface{}, string) {
	// 创建一个 HttpClient 实例，设置超时时间为 60 秒
	client := utils.NewHttpClient(60 * time.Second)
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

	var cateStr string
	var scripts []string
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
			cateStr = text[len("cgiData.appmsg_list =") : len(text)-len(".appmsg_list;")]
			break
		}
	}
	dataMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(cateStr), &dataMap)
	if err != nil {
		log.Println("cateStr：", cateStr)
		log.Println("解析响应结果错误：", err)
	}
	title := fmt.Sprintf("%s_%s", authorName, collectionName)
	return dataMap, title
}

func StartCollectionHome(urlStr string, path string) ([]Collect, error) {
	var collects []Collect
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
		collect := Collect{Name: jsName, StartPath: cateFilePath, EndPath: endFilePath}
		// 默认分类
		if i == 0 {
			// 页面容量
			begin += len(categoryList)
			collect.List = append(collect.List, categoryList...)
		}
		log.Println("开始下标：", begin)

		for {
			//time.Sleep(time.Second * 5)
			// 分页加载
			result := listPageHome(biz, strconv.Itoa(i), sn, sessionUs, hid, begin)
			if result != nil {
				articleList := result["appmsg_list"].([]interface{}) // 文章列表
				hasMore = int(result["has_more"].(float64))          // 结束标识 0结束 1继续
				for _, article := range articleList {
					articleMap, ok := article.(map[string]interface{})
					if !ok {
						panic("article is not a map")
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
		// 一次写入
		utils.WriteAppendFile(cateFilePath, strings.Join(utils.UrlFilter(collect.List), "\n"))
		fmt.Printf("分类：%s 写入完成\n", name)
	}
	return collects, nil
}
