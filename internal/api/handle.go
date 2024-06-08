package api

import (
	"encoding/json"
	"fmt"
	"go-wx-download/config"
	"go-wx-download/internal/biz"
	"go-wx-download/pkg/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type List struct {
	List []string `json:"list"`
}

func RemoveDuplicates(arr []string) []string {
	// 创建一个 map 来记录已经出现过的元素
	seen := make(map[string]struct{})
	// 创建一个切片来存储去重后的结果
	var result []string

	// 遍历原始数组
	for _, value := range arr {
		// 如果元素没有出现在 map 中，则添加到结果中
		if _, found := seen[value]; !found {
			seen[value] = struct{}{} // 标记元素已经出现过
			result = append(result, value)
		}
	}

	return result
}

// HandlerGather 采集接口
func HandlerGather(w http.ResponseWriter, r *http.Request, defaultDataPath string, port string, cfg *config.Config) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 设置响应头部为 JSON 类型
	w.Header().Set("Content-Type", "application/json")
	// 请求方式
	method := r.Method
	// 请求地址
	path := r.URL.Path
	// 请求参数
	query := r.URL.RawQuery
	// 采集地址列表
	var urlArray []string
	// 开始时间
	start := time.Now()

	log.Println("==============采集接口=======================")
	log.Printf("请求地址：%s\n", path)
	log.Printf("请求方式：%s\n", method)

	switch method {

	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("请求体读取异常！", err)
			_, _ = w.Write(utils.Fail("请求体读取异常"))
			return
		}
		err = json.Unmarshal(body, &urlArray)
		if err != nil {
			log.Println("请求体解析错误！", err)
			_, err = w.Write(utils.Fail("请求体解析错误！"))
			return
		}
		break
	default:
		// 删除 /gather/
		path = path[len("/gather/"):]
		// 修改 https://mp.weixin.qq.com 为 https://mp.weixin.qq.com
		path = strings.ReplaceAll(path, "https:/mp.weixin.qq.com", "https://mp.weixin.qq.com")
		// 拼接完整地址
		if len(query) > 0 {
			path += "?" + query
		}
		if len(path) > 0 {
			urlArray = append(urlArray, path)
		}
	}
	jsonData, err := json.Marshal(urlArray)
	if err != nil {
		log.Print("将映射封送至JSON时出错", jsonData)
	}
	log.Printf("请求参数/内容：%s,%s \n", query, jsonData)

	results := RemoveDuplicates(utils.UrlFilter(urlArray))
	if len(results) == 0 {
		_, _ = w.Write(utils.Fail("地址不符和规范！"))
		return
	}

	log.Printf("文章数量：%d\n", len(urlArray))
	fmt.Println("开始时间：", start)

	// 定义 goroutines 等待组
	var wg sync.WaitGroup
	// 并发数量
	maxConcurrency := utils.Iit(cfg.Thread.Html > 0, cfg.Thread.Html, 5)
	// 创建 sem 通道
	sem := make(chan struct{}, maxConcurrency)
	// 用于存储文件地址的通道
	filePaths := make(chan string, len(results))
	for i, item := range results {
		log.Printf("当前下标：[%d]\n", i+1)
		log.Printf("当前资源：[%s]\n", item)
		wg.Add(1)
		go biz.DownloadHtml(item, defaultDataPath, "", sem, &wg, filePaths, cfg.Thread.Image)
	}
	// 等待所有下载完成
	go func() {
		wg.Wait()
		close(filePaths)
	}()

	// 收集所有文件大小
	for f := range filePaths {
		log.Println(f)
		// 转PDF
		if cfg.Wkhtmltopdf.Enable && len(f) > 0 {
			if len(f) <= 0 {
				continue
			}
			fmt.Println("文件路径：" + f)
			split := strings.Split(f, string(os.PathSeparator))
			list := split[len(split)-3:]
			fmt.Println(list)
			// 第一个参数 端口； 第二个参数 公众号名称； 第三个参数 文件名称
			httpURL := fmt.Sprintf("http://127.0.0.1:%s/wx/%s/html/%s", port, url.QueryEscape(list[0]), url.QueryEscape(list[2]))
			f = filepath.Join(defaultDataPath, list[0], "pdf", list[2][0:len(list[2])-len(".html")]+".pdf")
			go utils.ToPDF(f, httpURL, cfg.Wkhtmltopdf.Path)
		}
	}

	// 结束时间
	end := time.Now()
	duration := end.Sub(start)
	fmt.Println("结束时间：", end)
	fmt.Println("采集耗时：", duration)
	// 返回 JSON 数据
	_, err = w.Write(utils.Success("结束"))
	if err != nil {
		return
	}
}

// HandlerFolder 目录下的文件
// folder 目录、fType 文件类型、path 路径
func HandlerFolder(path string, w http.ResponseWriter, r *http.Request) {
	folder := r.URL.Path[len("/ats/"):]
	folderType := ""
	if len(folder) > 0 {
		params := strings.Split(folder, "/")
		folder = params[0]
		if len(params) >= 2 {
			folderType = params[1]
		}
	}
	if len(folder) == 0 {
		folders := biz.GetFolders(path)
		data, err := json.Marshal(folders)
		_, err = w.Write([]byte(utils.Iif(len(folders) == 0, "[]", string(data))))
		if err != nil {
			return
		}
	} else if len(folderType) > 0 {
		folders := biz.GetFolderDetail(path, folder, folderType)
		data, err := json.Marshal(folders)
		_, err = w.Write([]byte(utils.Iif(len(folders) == 0, "[]", string(data))))
		if err != nil {
			return
		}
	}
}

func HandlerOpen(w http.ResponseWriter, r *http.Request, defaultDataPath string) {
	path := r.URL.Path
	pathSegments := strings.Split(path, "/")
	lastSegment := pathSegments[len(pathSegments)-1]
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		iif := utils.Iif(runtime.GOOS == "windows", "start", "open")
		utils.CmdOpenFolder(iif + " " + filepath.Join(defaultDataPath, lastSegment))
		_, err := w.Write(utils.Success("打开成功"))
		if err != nil {
			return
		}
	} else {
		log.Println("不支持 open")
		_, err := w.Write(utils.Success("不支持"))
		if err != nil {
			return
		}
	}
}

// HandlerCollect 采集接口
func HandlerCollect(w http.ResponseWriter, r *http.Request, defaultDataPath string, port string, cfg *config.Config) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 设置响应头部为 JSON 类型
	w.Header().Set("Content-Type", "application/json")
	// 请求方式
	method := r.Method
	// 请求地址
	path := r.URL.Path
	paramsMap := make(map[string]string)
	// 开始时间
	start := time.Now()

	log.Println("==============采集接口=======================")
	log.Printf("请求地址：%s\n", path)
	log.Printf("请求方式：%s\n", method)

	switch method {

	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("请求体读取异常！", err)
			_, _ = w.Write(utils.Fail("请求体读取异常"))
			return
		}
		err = json.Unmarshal(body, &paramsMap)
		if err != nil {
			log.Println("请求体解析错误！", err)
			_, err = w.Write(utils.Fail("请求体解析错误！"))
			return
		}
		break
	default:
		_, _ = w.Write(utils.Fail("请求方式错误！"))
		return
	}
	get := paramsMap["url"]
	folder := paramsMap["folder"]
	log.Printf("合集地址：%s\n", get)

	if !strings.Contains(get, "appmsgalbum") && !strings.Contains(get, "homepage") {
		_, _ = w.Write(utils.Fail("不支持该地址，貌似不是合集地址！"))
		return
	}
	var collects []biz.Collect
	if strings.Contains(get, "appmsgalbum") {
		album, err2 := biz.StartCollectionAlbum(get, defaultDataPath)
		if err2 != nil {
			_, _ = w.Write(utils.Fail("解析异常"))
			return
		}
		collects = append(collects, album)
	} else if strings.Contains(get, "homepage") {
		album, err2 := biz.StartCollectionHome(get, defaultDataPath)
		fmt.Println("home数量", len(album))
		if err2 != nil {
			_, _ = w.Write(utils.Fail("解析异常"))
			return
		}
		collects = append(collects, album...)
	}
	for i, album := range collects {
		urlArray := album.List
		results := RemoveDuplicates(utils.UrlFilter(urlArray))
		if len(results) == 0 {
			_, _ = w.Write(utils.Fail("地址不符和规范！"))
			fmt.Println("跳过：", i, urlArray)
			continue
		}

		log.Printf("文章数量：%d\n", len(urlArray))
		fmt.Println("开始时间：", start)

		// 定义 goroutines 等待组
		var wg sync.WaitGroup
		// 并发数量
		maxConcurrency := utils.Iit(cfg.Thread.Html > 0, cfg.Thread.Html, 5)
		// 创建 sem 通道
		sem := make(chan struct{}, maxConcurrency)
		// 用于存储文件地址的通道
		filePaths := make(chan string, len(results))
		for i, item := range results {
			log.Printf("当前下标：[%d]\n", i+1)
			log.Printf("当前资源：[%s]\n", item)
			wg.Add(1)
			go biz.DownloadHtml(item, defaultDataPath, utils.Iif(len(folder) > 0, folder, album.Name), sem, &wg, filePaths, cfg.Thread.Image)
		}
		// 等待所有下载完成
		go func() {
			wg.Wait()
			close(filePaths)
		}()

		// 收集所有文件大小
		for f := range filePaths {
			log.Println(f)
			// 转PDF
			if cfg.Wkhtmltopdf.Enable && len(f) > 0 {
				if len(f) <= 0 {
					continue
				}
				fmt.Println("文件路径：" + f)
				split := strings.Split(f, string(os.PathSeparator))
				list := split[len(split)-3:]
				fmt.Println(list)
				// 第一个参数 端口； 第二个参数 公众号名称； 第三个参数 文件名称
				httpURL := fmt.Sprintf("http://127.0.0.1:%s/wx/%s/html/%s", port, url.QueryEscape(list[0]), url.QueryEscape(list[2]))
				f = filepath.Join(defaultDataPath, list[0], "pdf", list[2][0:len(list[2])-len(".html")]+".pdf")
				go utils.ToPDF(f, httpURL, cfg.Wkhtmltopdf.Path)
			}
		}

		// 结束时间
		end := time.Now()
		duration := end.Sub(start)
		fmt.Println("结束时间：", end)
		fmt.Println("采集耗时：", duration)
		// 写新文件
		utils.CopyFile(album.EndPath, album.StartPath)
		// 删除
		os.Remove(album.StartPath)
	}
	// 返回 JSON 数据
	_, err := w.Write(utils.Success("结束"))
	if err != nil {
		return
	}
}
