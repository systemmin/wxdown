package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// DownloadFile 下载文件多线程
func DownloadFile(url string, filepath string, headers map[string]string, sem chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	// 从信号量中获取一个令牌
	sem <- struct{}{}

	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("无法创建文件：%s\n", err)
		// 释放信号量令牌
		<-sem
		return
	}
	defer file.Close()

	// 创建请求客户端对象
	client := &http.Client{
		Timeout: 5 * 60 * time.Second, // 设置超时时间为10分钟
	}

	// 创建一个请求实例
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求实例时出错:", err)
		return
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	// 设置请求头
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
	// 下载文件
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("无法下载文件：%s\n", err)
		// 释放信号量令牌
		<-sem
		return
	}
	defer response.Body.Close()

	// 获取文件大小
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Printf("无法写入文件：%s\n", err)
		// 释放信号量令牌
		<-sem
		return
	}
	log.Println(filepath)

	// 释放信号量令牌
	<-sem
}
