package down

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"go-wx-download/pkg/utils"
	"golang.org/x/image/webp"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// DownloadFile 多线程下载文件
func DownloadFile(url string, filepath string, headers map[string]string, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// 从信号量中获取一个令牌
	sem <- struct{}{}
	defer func() { <-sem }() // 确保在函数返回时释放信号量令牌

	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("无法创建文件：%s\n", err)
		return
	}
	defer file.Close()

	// 创建一个 HttpClient 实例，设置超时时间为 10 分钟
	client := utils.NewHttpClient(10 * 60 * time.Second)
	// 发送请求
	response, err := client.Request("GET", url, headers, nil)
	if err != nil {
		fmt.Println("创建请求实例时出错:", err)
		return
	}

	// 拷贝文件到本地
	_, err = io.Copy(file, bytes.NewReader(response)) // 字节转 reader 对象
	if err != nil {
		fmt.Printf("无法写入文件：%s\n", err)
		return
	}
	webPToJPEG(filepath)
	log.Println(filepath)
}

// ImageToBase64 图片转base64
func ImageToBase64(url, suffix string, headers map[string]string, sem chan struct{}, wg *sync.WaitGroup) string {
	defer wg.Done()

	// 从信号量中获取一个令牌
	sem <- struct{}{}
	defer func() { <-sem }() // 确保在函数返回时释放信号量令牌
	// 创建一个 HttpClient 实例，设置超时时间为 10 分钟
	client := utils.NewHttpClient(10 * 60 * time.Second)
	// 发送请求
	response, err := client.Request("GET", url, headers, nil)
	if err != nil {
		fmt.Println("创建请求实例时出错:", err)
		return ""
	}
	// 将图片内容转换为Base64编码
	base64Str := base64.StdEncoding.EncodeToString(response)

	// 拼接前缀
	base64WithPrefix := "data:image/" + suffix + ";base64," + base64Str

	return base64WithPrefix
}

// webpToJPG webp 转 jpeg 格式
func webPToJPEG(filepath string) {
	before, b := strings.CutSuffix(filepath, "webp")
	if b {
		// 打开WebP文件
		webpFile, err := os.Open(filepath)
		if err != nil {
			fmt.Println("打开文件错误:", err)
			return
		}
		defer webpFile.Close()
		// 解码 WebP 文件为图像对象
		img, err := webp.Decode(webpFile)
		if err != nil {
			fmt.Println("图片解码错误:", err)
			return
		}

		// 创建一个新的JPG文件
		jpgFile, err := os.Create(fmt.Sprintf("%sjpeg", before))
		if err != nil {
			fmt.Println("创建图错误:", err)
			return
		}
		defer jpgFile.Close()

		// 将图像对象编码为JPG格式并写入文件
		err = imaging.Encode(jpgFile, img, imaging.JPEG)
		if err != nil {
			fmt.Println("图片编码错误:", err)
			return
		}
		fmt.Println("格式转换成功!")
	}
}
