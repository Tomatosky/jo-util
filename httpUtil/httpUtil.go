package httpUtil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Tomatosky/jo-util/strUtil"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

type Resp struct {
	Err        error
	Body       []byte
	StatusCode int
	Headers    map[string]string
}

func (r *Resp) Raw() []byte {
	return r.Body
}

func (r *Resp) Text() string {
	return string(r.Body)
}

func (r *Resp) Json() (map[string]any, error) {
	jsonData := map[string]any{}
	err := json.Unmarshal(r.Body, &jsonData)
	return jsonData, err
}

type GetOptions struct {
	Headers map[string]string
	Timeout int
}

// Get 发送GET请求
func (rc *RequestClient) Get(url string, getOptions *GetOptions) *Resp {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Resp{Err: err}
	}

	// 设置请求头
	if getOptions != nil && getOptions.Headers != nil {
		for key, value := range getOptions.Headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求
	if getOptions != nil && getOptions.Timeout > 0 {
		rc.Client.Timeout = time.Duration(getOptions.Timeout) * time.Second
	}
	resp, err := rc.Client.Do(req)
	if err != nil {
		return &Resp{Err: err}
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Resp{Err: err}
	}

	// 收集响应头
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0] // 只取第一个值
		}
	}

	return &Resp{
		Err:        nil,
		Body:       body,
		StatusCode: resp.StatusCode,
		Headers:    headers,
	}
}

type PostOptions struct {
	Headers     map[string]string
	Timeout     int
	IsJson      bool
	IsMultipart bool
}

// Post 发送POST请求
func (rc *RequestClient) Post(postUrl string, data map[string]interface{}, postOptions *PostOptions) *Resp {
	var (
		reqBody     io.Reader
		contentType string
	)
	// 判断请求体格式
	switch {
	case postOptions.IsMultipart:
		// multipart/form-data 格式
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for key, value := range data {
			// 处理文件上传情况
			if file, ok := value.(*os.File); ok {
				part, err := writer.CreateFormFile(key, filepath.Base(file.Name()))
				if err != nil {
					return &Resp{Err: err}
				}
				_, err = io.Copy(part, file)
				if err != nil {
					return &Resp{Err: err}
				}
				continue
			}

			// 处理普通字段
			strValue := fmt.Sprintf("%v", value)
			err := writer.WriteField(key, strValue)
			if err != nil {
				return &Resp{Err: err}
			}
		}

		err := writer.Close()
		if err != nil {
			return &Resp{Err: err}
		}

		reqBody = body
		contentType = writer.FormDataContentType()

	case postOptions.IsJson:
		// JSON 格式
		jsonData, err := json.Marshal(data)
		if err != nil {
			return &Resp{Err: err}
		}
		reqBody = bytes.NewBuffer(jsonData)
		contentType = "application/json"

	default:
		// 表单格式 (x-www-form-urlencoded)
		formData := url.Values{}
		for key, value := range data {
			formData.Set(key, fmt.Sprintf("%v", value)) // 确保值转为字符串
		}
		reqBody = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", postUrl, reqBody)
	if err != nil {
		return &Resp{Err: err}
	}
	// 设置请求头
	req.Header.Set("Content-Type", contentType)
	if postOptions.Headers != nil {
		for key, value := range postOptions.Headers {
			req.Header.Set(key, value)
		}
	}
	// 设置超时（如果配置）
	if postOptions.Timeout > 0 {
		rc.Client.Timeout = time.Duration(postOptions.Timeout) * time.Second
	}
	// 发送请求
	resp, err := rc.Client.Do(req)
	if err != nil {
		return &Resp{Err: err}
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Resp{Err: err}
	}

	// 收集响应头
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0] // 只取第一个值
		}
	}

	return &Resp{
		Err:        nil,
		Body:       body,
		StatusCode: resp.StatusCode,
		Headers:    headers,
	}
}

func UrlEncode(params map[string]interface{}) string {
	var p = url.Values{}
	for k, v := range params {
		p.Add(k, strUtil.ToString(v))
	}
	return p.Encode()
}
