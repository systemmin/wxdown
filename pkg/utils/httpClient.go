package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// HttpClient 封装了 HTTP 客户端
type HttpClient struct {
	client *http.Client
}

// NewHttpClient 创建一个新的 HttpClient 实例
func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Request 执行一个 HTTP 请求
func (hc *HttpClient) Request(method, url string, headers map[string]string, body []byte) ([]byte, error) {
	log.Println("====================http====================")
	// 创建一个新的 HTTP 请求
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	// 添加请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	log.Println("====================响应header====================")
	log.Println("请求地址：", url)
	log.Println("请求方式：", method)
	// 执行请求
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		req.Header.Set(key, value[0])
		log.Printf("%s:%s", key, value)
	}
	log.Println("响应状态：", resp.StatusCode)

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
