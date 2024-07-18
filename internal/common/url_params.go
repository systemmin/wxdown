package common

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type UrlParams struct {
	Folder string   `json:"folder"`
	Urls   []string `json:"urls"`
}

func ReadBody(b io.ReadCloser, urlParams *UrlParams) error {
	body, err := io.ReadAll(b)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &urlParams)
	if err != nil {
		return err
	}
	return nil
}

func GetParams(r *http.Request) (UrlParams, error) {
	query := r.URL.RawQuery
	path := r.URL.Path
	var urlParams UrlParams
	switch r.Method {
	case "POST":
		err := ReadBody(r.Body, &urlParams)
		if err != nil {
			return urlParams, err
		}
	default:
		index := strings.Index(path, "http")
		if index != -1 {
			path = path[index:]
		}
		path = strings.ReplaceAll(path, "https:/mp.weixin.qq.com", "https://mp.weixin.qq.com")
		// 拼接完整地址
		if len(query) > 0 {
			path += "?" + query
		}
		if len(path) > 0 {
			urlParams.Urls = append(urlParams.Urls, path)
		}
	}
	return urlParams, nil
}
