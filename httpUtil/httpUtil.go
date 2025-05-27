package httpUtil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Tomatosky/jo-util/strUtil"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	var (
		reqBody     io.Reader
		contentType string
	)

	// 判断是否使用 JSON 格式
	if postOptions.IsJson {
		// JSON 格式
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		reqBody = bytes.NewBuffer(jsonData)
		contentType = "application/json"
	} else {
		// 表单格式 (x-www-form-urlencoded)
		formData := url.Values{}
		for key, value := range data {
			formData.Set(key, fmt.Sprintf("%v", value)) // 确保值转为字符串
		}
		reqBody = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return "", err
	}
	// 设置请求头
	if postOptions.Headers != nil {
		for key, value := range postOptions.Headers {
			req.Header.Set(key, value)
		}
	}
	req.Header.Set("Content-Type", contentType) // 根据 IsJson 设置正确的 Content-Type
	// 设置超时（如果配置）
	if postOptions.Timeout > 0 {
		rc.Client.Timeout = time.Duration(postOptions.Timeout) * time.Second
	}
	// 发送请求
	resp, err := rc.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应
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
