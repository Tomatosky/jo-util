package httpUtil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/Tomatosky/jo-util/strUtil"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RequestClient 是一个HTTP客户端工具类
type RequestClient struct {
	Client *http.Client
}

// NewRequestClient 创建一个新的RequestClient实例
func NewRequestClient() *RequestClient {
	return &RequestClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

type GetOptions struct {
	Headers map[string]string
	Timeout int
}

// Download 发送GET请求
func (rc *RequestClient) Download(url string, getOptions *GetOptions) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	if getOptions.Headers != nil {
		for key, value := range getOptions.Headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求
	if getOptions.Timeout > 0 {
		rc.Client.Timeout = time.Duration(getOptions.Timeout) * time.Second
	}
	resp, err := rc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Get 发送GET请求
func (rc *RequestClient) Get(url string, getOptions *GetOptions) (string, error) {
	body, err := rc.Download(url, getOptions)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type PostOptions struct {
	Headers map[string]string
	Timeout int
	IsJson  bool
}

// Post 发送POST请求
func (rc *RequestClient) Post(url string, data map[string]interface{}, postOptions *PostOptions) (string, error) {
	// 将数据转换为JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// 设置请求头
	if postOptions.Headers != nil {
		for key, value := range postOptions.Headers {
			req.Header.Set(key, value)
		}
	}
	if postOptions.IsJson {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	if postOptions.Timeout > 0 {
		rc.Client.Timeout = time.Duration(postOptions.Timeout) * time.Second
	}
	resp, err := rc.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func UrlEncode(params map[string]interface{}) string {
	var p = url.Values{}
	for k, v := range params {
		p.Add(k, strUtil.ToString(v))
	}
	return p.Encode()
}
