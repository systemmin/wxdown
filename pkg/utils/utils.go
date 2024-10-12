package utils

import (
	"fmt"
	"go-wx-download/internal/constant"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// SanitizeFilename 处理名称特殊字符
func SanitizeFilename(filename string) string {
	// 使用正则表达式替换掉不允许的字符为空字符
	sanitizedFilename := regexp.MustCompile(`[<>:"/\\|?*\r\n\s\*…]`).ReplaceAllString(filename, "")

	// 在Windows系统中，还应该移除 \ 和 /，因为它们是路径分隔符
	sanitizedFilename = regexp.MustCompile(`[\\/]`).ReplaceAllString(sanitizedFilename, "")

	// 移除表情符号
	sanitizedFilename = regexp.MustCompile("[\U0001F600-\U0001F64F\U0001F300-\U0001F5FF\U0001F680-\U0001F6FF\u2600-\u26FF\u2700-\u27BF\U0001F900-\U0001F9FF\U0001FA70-\U0001FAFF]").ReplaceAllString(sanitizedFilename, "")

	return sanitizedFilename
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

// ToPDF html 转 PDF
// path 输出路径
// url 访问地址
// wk bin 路径
func ToPDF(path string, url string, wk string) {
	// 改变工作目录
	if len(wk) > 0 {
		err := os.Chdir(wk)
		if err != nil {
			fmt.Println("改变目录失败:", err)
			return
		}
	}
	if runtime.GOOS != "windows" && strings.Contains(path, " ") {
		path = fmt.Sprintf(`"%s"`, path)
	}
	ExecuteCmd("wkhtmltopdf", url, path)
}

// ExecuteCmd 执行 cmd 命令
func ExecuteCmd(cmd string, args ...string) {
	fmt.Println("cmd:", cmd, "args:", args)
	command := exec.Command(cmd, args...)
	// 执行命令
	out, err := command.CombinedOutput() // 获取输出和错误信息
	if err != nil {
		fmt.Printf("执行命令时出错: %v\n输出: %s\n", err, string(out))
		return
	}
	// 输出命令执行结果
	fmt.Println("命令执行结果：", string(out))
}

func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
	}
}

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

func IsURL(url string) error {
	client := http.Client{}
	if _, err := client.Get(url); err != nil {
		return err
	}
	return nil
}
