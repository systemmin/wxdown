package down

import (
	"bytes"
	"fmt"
	"go-wx-download/pkg/utils"
	"io"
	"log"
	"os"
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
	log.Println(filepath)
}
