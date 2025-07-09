package httpUtil

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	// 创建一个测试HTTP服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法
		if r.Method != "GET" {
			t.Error("Expected GET request")
		}

		// 测试不同的路径返回不同的响应
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success response"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error response"))
		case "/headers":
			// 测试请求头
			if r.Header.Get("Test-Header") != "test-value" {
				t.Error("Expected Test-Header to be 'test-value'")
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("headers test"))
		case "/timeout":
			// 测试超时
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("should timeout"))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		}
	}))
	defer ts.Close()

	// 创建测试客户端
	client := NewRequestClient()

	tests := []struct {
		name     string
		url      string
		options  *GetOptions
		wantErr  bool
		wantBody string
	}{
		{name: "成功请求", url: ts.URL + "/success", options: nil, wantErr: false, wantBody: "success response"},
		{name: "服务器错误", url: ts.URL + "/error", options: nil, wantErr: false, wantBody: "error response"},
		{name: "带请求头的请求", url: ts.URL + "/headers", options: &GetOptions{Headers: map[string]string{"Test-Header": "test-value"}}, wantErr: false, wantBody: "headers test"},
		{name: "超时请求", url: ts.URL + "/timeout", options: &GetOptions{Timeout: 1}, wantErr: true},
		{name: "无效URL", url: "http://invalid.url", options: nil, wantErr: true, wantBody: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := client.Get(tt.url, tt.options)

			if tt.wantErr {
				if resp.Err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if resp.Err != nil {
				t.Errorf("Unexpected error: %v", resp.Err)
				return
			}

			if got := resp.Text(); got != tt.wantBody {
				t.Errorf("Expected body '%s', got '%s'", tt.wantBody, got)
			}
		})
	}
}

func TestRespMethods(t *testing.T) {
	// 测试Resp的各种方法
	jsonBody := `{"key":"value","num":123}`
	resp := &Resp{
		Body: []byte(jsonBody),
	}

	t.Run("Text", func(t *testing.T) {
		if resp.Text() != jsonBody {
			t.Error("Text() did not return expected body")
		}
	})

	t.Run("Json", func(t *testing.T) {
		jsonData, _ := resp.Json()
		if jsonData["key"] != "value" || jsonData["num"] != float64(123) {
			t.Error("Json() did not return expected data")
		}
	})

	t.Run("Json with invalid data", func(t *testing.T) {
		res := &Resp{
			Body: []byte("invalid json"),
		}
		// 这里不会panic，但返回的map可能是空的
		_, _ = res.Json()
	})
}

func TestPost(t *testing.T) {
	// 创建一个测试用的HTTP服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查请求头
		contentType := r.Header.Get("Content-Type")

		// 读取请求体
		_, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error("读取请求体失败:", err)
			return
		}

		// 根据不同的Content-Type返回不同的响应
		switch {
		case strings.Contains(contentType, "multipart/form-data"):
			w.Write([]byte("multipart response"))
		case strings.Contains(contentType, "application/json"):
			w.Write([]byte("json response"))
		case strings.Contains(contentType, "application/x-www-form-urlencoded"):
			w.Write([]byte("form response"))
		default:
			w.Write([]byte("unknown format"))
		}
	}))
	defer ts.Close()

	// 创建测试用的临时文件
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal("创建临时文件失败:", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test file content")
	tmpFile.Close()

	// 测试用例
	tests := []struct {
		name    string
		data    map[string]interface{}
		options *PostOptions
		want    string
		wantErr bool
	}{
		{
			name: "测试JSON格式请求",
			data: map[string]interface{}{"key1": "value1", "key2": 2},
			options: &PostOptions{
				IsJson: true,
				Headers: map[string]string{
					"X-Custom-Header": "test",
				},
			},
			want:    "json response",
			wantErr: false,
		},
		{
			name: "测试表单格式请求",
			data: map[string]interface{}{"key1": "value1", "key2": 2},
			options: &PostOptions{
				IsJson:      false,
				IsMultipart: false,
			},
			want:    "form response",
			wantErr: false,
		},
		{
			name: "测试Multipart格式请求(带文件)",
			data: map[string]interface{}{
				"file": func() *os.File {
					f, _ := os.Open(tmpFile.Name())
					return f
				}(),
				"field": "value",
			},
			options: &PostOptions{
				IsMultipart: true,
			},
			want:    "multipart response",
			wantErr: false,
		},
		{
			name: "测试Multipart格式请求(不带文件)",
			data: map[string]interface{}{"key1": "value1", "key2": 2},
			options: &PostOptions{
				IsMultipart: true,
			},
			want:    "multipart response",
			wantErr: false,
		},
		{
			name: "测试超时设置",
			data: map[string]interface{}{"key": "value"},
			options: &PostOptions{
				Timeout: 1, // 1秒超时
			},
			want:    "form response",
			wantErr: false,
		},
		{
			name:    "测试无效URL",
			data:    map[string]interface{}{"key": "value"},
			options: &PostOptions{},
			want:    "",
			wantErr: true,
		},
	}

	rc := NewRequestClient()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := ts.URL
			if tt.name == "测试无效URL" {
				url = "http://invalid.url"
			}

			resp := rc.Post(url, tt.data, tt.options)

			if (resp.Err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", resp.Err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp.Text() != tt.want {
				t.Errorf("Post() = %v, want %v", resp.Text(), tt.want)
			}
		})
	}
}

func TestPost_FileHandling(t *testing.T) {
	// 创建一个测试用的HTTP服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析multipart表单
		err := r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			t.Error("解析multipart表单失败:", err)
			return
		}

		// 检查文件是否存在
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Error("获取上传文件失败:", err)
			return
		}
		defer file.Close()

		// 检查文件名
		if header.Filename != filepath.Base(header.Filename) {
			t.Error("文件名不正确")
		}

		// 检查文件内容
		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		if buf.String() != "test file content" {
			t.Error("文件内容不正确")
		}

		w.Write([]byte("file upload success"))
	}))
	defer ts.Close()

	// 创建测试用的临时文件
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal("创建临时文件失败:", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test file content")
	tmpFile.Close()

	// 打开文件用于测试
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatal("打开临时文件失败:", err)
	}
	defer file.Close()

	rc := NewRequestClient()
	resp := rc.Post(ts.URL, map[string]interface{}{
		"file":  file,
		"field": "value",
	}, &PostOptions{
		IsMultipart: true,
	})

	if resp.Err != nil {
		t.Error("文件上传请求失败:", resp.Err)
	}

	if resp.Text() != "file upload success" {
		t.Error("文件上传响应不正确:", resp.Text())
	}
}

func TestPost_Timeout(t *testing.T) {
	// 创建一个慢速响应的测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 延迟2秒响应
		w.Write([]byte("response"))
	}))
	defer ts.Close()

	rc := NewRequestClient()

	// 设置1秒超时
	resp := rc.Post(ts.URL, map[string]interface{}{"key": "value"}, &PostOptions{
		Timeout: 1,
	})

	if resp.Err == nil {
		t.Error("预期超时错误，但未发生")
	}
}

func TestUrlEncode(t *testing.T) {
	params := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	encoded := UrlEncode(params)
	expected := "key1=value1&key2=123&key3=true"
	if encoded != expected {
		t.Errorf("Expected '%s', got '%s'", expected, encoded)
	}
}
