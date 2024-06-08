package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-wx-download/internal/constant"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// SanitizeFilename 处理名称特殊字符
func SanitizeFilename(filename string) string {
	// 使用正则表达式替换掉不允许的字符为空字符
	sanitizedFilename := regexp.MustCompile(`[<>:"/\\|?*\r\n\s\*…]`).ReplaceAllString(filename, "")

	// 在Windows系统中，还应该移除 \ 和 /，因为它们是路径分隔符
	sanitizedFilename = regexp.MustCompile(`[\\/]`).ReplaceAllString(sanitizedFilename, "")

	return sanitizedFilename
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

// UrlFilter 过滤采集地址
func UrlFilter(urls []string) []string {

	var res []string
	for _, url := range urls {
		if len(url) > 0 && strings.Contains(url, constant.Domain) {
			// 处理特殊字符
			if strings.Contains(url, "&amp;") {
				res = append(res, strings.ReplaceAll(url, "&amp;", "&"))
			} else {
				res = append(res, url)
			}
		}
	}
	return res
}

// CreateNewFolder 在 data 目录下创建新文件夹
func CreateNewFolder(folderName string, path string) {
	// 创建新文件夹
	joinPath := filepath.Join(path, folderName)
	if _, err := os.Stat(joinPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("目录 '%s' 不存在\n", joinPath)
		} else {
			fmt.Printf("无法检查目录 '%s': %v\n", joinPath, err)
		}
	} else {
		fmt.Printf("目录 '%s' 存在\n", joinPath)
		return
	}

	os.Mkdir(joinPath, os.ModePerm)
	// 数据 文件夹
	os.Mkdir(filepath.Join(joinPath, "db"), os.ModePerm)
	// 图片 文件夹
	os.Mkdir(filepath.Join(joinPath, "images"), os.ModePerm)
	// pdf 文件夹
	os.Mkdir(filepath.Join(joinPath, "pdf"), os.ModePerm)
	// html 文件夹
	os.Mkdir(filepath.Join(joinPath, "html"), os.ModePerm)
	// audio 音频
	os.Mkdir(filepath.Join(joinPath, "audios"), os.ModePerm)
	// video 视频
	os.Mkdir(filepath.Join(joinPath, "videos"), os.ModePerm)

}

// CopyFile dst 目标文件, src 源文件
func CopyFile(dst string, src string) {
	// 1、判断目标文件是否存在
	file, _ := os.Stat(dst)
	if file != nil {
		//log.Printf("目标文件 '%s' 已存在\n", dst)
		return
	}
	// 2、创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Printf("创建文件 '%s' 失败; 描述信息：%s\n", dst, err)
	}
	// 3、打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		log.Printf("源文件 '%s' 不存在; 描述信息：%s\n", dst, err)
		return
	}
	// 函数返回前正确关闭
	defer srcFile.Close()
	// 4、拷贝
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Printf("拷贝文件 '%s' 到 '%s' 失败; 描述信息：%s\n", src, dst, err)
	}
}

// WriteAppendFile 追加写入
func WriteAppendFile(path, content string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()
	// 写入内容
	if _, err := file.WriteString(content); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}
}

// GetBgImage 获取 css 背景图片
func GetBgImage(style string) string {
	str := style[strings.Index(style, "(")+1 : strings.LastIndex(style, ")")]
	if strings.HasPrefix(str, "\"") {
		str = str[1:]
		if strings.HasSuffix(str, "\"") {
			str = str[:len(str)-1]
		}
	}
	return str
}

// GetUrlParams 获取 URL 参数
func GetUrlParams(urlStr string, key string) string {
	// 使用 net/url 包中的 Parse 函数解析 URL
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("URL 解析错误:", err)
		return ""
	}

	// 使用 ParseQuery 函数解析查询字符串
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		fmt.Println("Error parsing the query:", err)
		return ""
	}
	return Iif(q.Get(key) == "other", q.Get("tp"), q.Get(key))
}

func Iif(b bool, get string, get2 string) string {
	if b {
		return get
	} else {
		return get2
	}
}

func Iit(b bool, get int, get2 int) int {
	if b {
		return get
	} else {
		return get2
	}
}

// GetSuffix 获取图片后缀
func GetSuffix(urlStr string) string {
	split := strings.Split(urlStr, ",")
	if len(split) >= 2 {
		return strings.Split(split[2], "_")[1]
	} else {
		return "png"
	}
}

// IsNotExistCreate 不存在就创建目录
func IsNotExistCreate(dirPath string) {
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); err != nil {
		// 如果目录不存在，创建它
		if os.IsNotExist(err) {
			// 创建目录
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				return
			}
			fmt.Printf("目录 %s 创建成功.\n", dirPath)
		} else {
			// 如果是因为其他原因导致错误，比如权限问题
			fmt.Printf("检查目录是否存在时出错: %v\n", err)
			return
		}
	}
}

func Upgradation(path string, v string) {
	fmt.Printf("版本：%s 升级,自动检测缺失文件夹,可忽略\n", v)
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Printf("%s 路径不存在\n", path)
	}
	for _, file := range dir {
		if file.IsDir() {
			name := file.Name()
			if name == ".DS_Store" || name == "css" || name == "task" {
				continue
			}
			IsNotExistCreate(filepath.Join(path, name, "audios"))
			IsNotExistCreate(filepath.Join(path, name, "videos"))
		}
	}
}

// ToPDF html 转 PDF
// cmd.exe /c cd C:\\your\\directory && your-command
// cmd.exe /k
// "sh", "-c", "cd /path/to/directory && your-command"
func ToPDF(path string, url string, wk string) {

	// 待执行命令
	env := Iif(runtime.GOOS == "windows", "cmd", "sh")
	quit := Iif(runtime.GOOS == "windows", "/c", "-c")
	suffix := Iif(len(wk) > 0, fmt.Sprintf("cd %s && wkhtmltopdf", wk), "wkhtmltopdf")

	cmd := fmt.Sprintf("%s %s %s", suffix, url, path)
	fmt.Println("env：", env)
	fmt.Println("quit：", quit)
	fmt.Println("html to pdf cmd：", cmd)
	// E:\Program Files\wkhtmltopdf\bin
	command := exec.Command(env, quit, cmd)
	// 创建一个缓冲区，用于存储命令的标准输出
	var out bytes.Buffer
	command.Stdout = &out
	// 执行命令
	err := command.Run()
	if err != nil {
		return
	}
	// 输出命令执行结果
	fmt.Println("命令执行结果：", out.String())
}

func CmdOpenFolder(cmd string) {
	fmt.Println("open folder cmd：", cmd)
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd", "/c", cmd)
	} else {
		command = exec.Command("sh", "-c", cmd)
	}
	// 创建一个缓冲区，用于存储命令的标准输出
	var out bytes.Buffer
	command.Stdout = &out
	// 执行命令
	err := command.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	// 输出命令执行结果
	fmt.Println("命令执行结果：", out.String())
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
	Ok      bool   `json:"ok"`
}

// Ok 返回状态码
func Ok(msg string) []byte {
	result := Result{
		Code:    200,
		Message: Iif(len(msg) > 0, msg, "成功"),
		Data:    "",
		Ok:      true,
	}
	jsonData, _ := json.Marshal(result)
	return jsonData
}

func Success(msg string) []byte {
	result := Result{
		Code:    200,
		Message: Iif(len(msg) > 0, msg, "成功"),
		Data:    "",
		Ok:      true,
	}
	jsonData, _ := json.Marshal(result)
	return jsonData
}

// Fail 返回状态码
func Fail(msg string) []byte {
	result := Result{
		Code:    500,
		Message: Iif(len(msg) > 0, msg, "失败"),
		Data:    "",
		Ok:      false,
	}
	jsonData, _ := json.Marshal(result)
	return jsonData
}

// ParseUrl 解析 URL 参数 返回 Params 对象
func ParseUrl(urlStr string) (Params, error) {
	params := make(Params)
	result, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	for key, values := range result.Query() {
		params.Set(key, values[0])
	}
	return params, nil
}
