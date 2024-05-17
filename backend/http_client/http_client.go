package http_client

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// TODO: 支持使用系统代理
var httpClientInst *http.Client
var once sync.Once

func initHttpClient() {
	httpClientInst = &http.Client{}
}

// HttpClientInst 返回httpClient的单例实例
func HttpClientInst() *http.Client {
	once.Do(initHttpClient)
	return httpClientInst
}

// UpdateProxy 设置httpClient的代理地址
func UpdateProxy(proxyUrl string) error {
	if proxyUrl == "" {
		return nil
	}
	pUrl, err := url.Parse(proxyUrl)
	if err != nil {
		return fmt.Errorf("parse proxy url failed: %w", err)
	}

	httpClient := HttpClientInst()
	httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(pUrl)}
	return nil
}
