package utils

import "strings"

type Params map[string]string

func (p Params) Set(key string, value string) {
	p[key] = value
}

func (p Params) Get(key string) string {
	return p[key]
}

func (p Params) ToString(url string) string {
	var params []string
	for k, v := range p {
		params = append(params, k+"="+v)
	}
	if len(url) > 0 {
		return url + "?" + strings.Join(params, "&")
	}
	return "?" + strings.Join(params, "&")
}
