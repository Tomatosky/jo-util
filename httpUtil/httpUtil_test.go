package httpUtil

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewRequestClient(t *testing.T) {
	client := NewRequestClient()
	if client.Client == nil {
		t.Error("Expected http.Client to be initialized, got nil")
	}
}

func TestDownload(t *testing.T) {
	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}))
	defer ts.Close()

	client := NewRequestClient()

	// 测试正常下载
	body, err := client.Download(ts.URL, &GetOptions{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if string(body) != "test response" {
		t.Errorf("Expected 'test response', got '%s'", string(body))
	}

	// 测试带header的下载
	headers := map[string]string{"X-Test": "value"}
	body, err = client.Download(ts.URL, &GetOptions{Headers: headers})
	if err != nil {
		t.Errorf("Expected no error with headers, got %v", err)
	}

	// 测试超时
	_, err = client.Download(ts.URL, &GetOptions{Timeout: 1})
	if err != nil {
		t.Errorf("Expected no error with timeout, got %v", err)
	}

	// 测试无效URL
	_, err = client.Download("invalid-url", &GetOptions{})
	if err == nil {
		t.Error("Expected error with invalid URL, got nil")
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test get response"))
	}))
	defer ts.Close()

	client := NewRequestClient()

	// 测试正常GET请求
	response, err := client.Get(ts.URL, &GetOptions{})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response != "test get response" {
		t.Errorf("Expected 'test get response', got '%s'", response)
	}
}

func TestPost(t *testing.T) {
	// 创建一个测试用的HTTP服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查Content-Type头
		contentType := r.Header.Get("Content-Type")

		// 根据不同的Content-Type处理请求
		switch {
		case strings.Contains(contentType, "multipart/form-data"):
			// 正确解析multipart表单的方法
			reader, err := r.MultipartReader()
			if err != nil {
				t.Error("创建MultipartReader失败:", err)
				return
			}

			// 读取multipart内容
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Error("读取part失败:", err)
					return
				}

				// 这里可以检查part的内容
				_, err = io.ReadAll(part)
				if err != nil {
					t.Error("读取part内容失败:", err)
					return
				}
			}
			w.Write([]byte("multipart received"))

		case contentType == "application/json":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error("读取请求体失败:", err)
				return
			}
			w.Write([]byte("json received: " + string(body)))

		case contentType == "application/x-www-form-urlencoded":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error("读取请求体失败:", err)
				return
			}
			w.Write([]byte("form received: " + string(body)))
		}
	}))
	defer ts.Close()

	// 创建测试用的RequestClient
	rc := NewRequestClient()

	// 测试用例1: 普通表单POST请求
	t.Run("普通表单POST", func(t *testing.T) {
		data := map[string]interface{}{
			"name":  "张三",
			"age":   25,
			"email": "zhangsan@example.com",
		}
		options := &PostOptions{
			Headers: map[string]string{
				"X-Test-Header": "test-value",
			},
			Timeout: 5,
		}

		resp, err := rc.Post(ts.URL, data, options)
		if err != nil {
			t.Error("POST请求失败:", err)
			return
		}

		if !strings.Contains(resp, "form received") {
			t.Error("响应不符合预期:", resp)
		}
	})

	// 测试用例2: JSON POST请求
	t.Run("JSON POST", func(t *testing.T) {
		data := map[string]interface{}{
			"name":  "李四",
			"age":   30,
			"email": "lisi@example.com",
		}
		options := &PostOptions{
			IsJson: true,
		}

		resp, err := rc.Post(ts.URL, data, options)
		if err != nil {
			t.Error("POST请求失败:", err)
			return
		}

		if !strings.Contains(resp, "json received") {
			t.Error("响应不符合预期:", resp)
		}
	})

	// 测试用例3: Multipart POST请求(包含文件上传)
	t.Run("Multipart POST with file", func(t *testing.T) {
		// 创建一个临时文件用于测试
		tempFile, err := os.CreateTemp("", "testfile-*.txt")
		if err != nil {
			t.Fatal("创建临时文件失败:", err)
		}
		defer os.Remove(tempFile.Name())

		_, err = tempFile.WriteString("这是一个测试文件内容")
		if err != nil {
			t.Fatal("写入临时文件失败:", err)
		}
		tempFile.Close()

		// 重新打开文件用于上传
		file, err := os.Open(tempFile.Name())
		if err != nil {
			t.Fatal("打开临时文件失败:", err)
		}
		defer file.Close()

		data := map[string]interface{}{
			"username": "王五",
			"avatar":   file,
		}
		options := &PostOptions{
			IsMultipart: true,
		}

		resp, err := rc.Post(ts.URL, data, options)
		if err != nil {
			t.Error("POST请求失败:", err)
			return
		}

		if resp != "multipart received" {
			t.Error("响应不符合预期:", resp)
		}
	})

	// 测试用例4: 超时测试
	t.Run("超时测试", func(t *testing.T) {
		// 创建一个会延迟响应的测试服务器
		slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second) // 延迟2秒响应
			w.Write([]byte("ok"))
		}))
		defer slowServer.Close()

		data := map[string]interface{}{
			"test": "timeout",
		}
		options := &PostOptions{
			Timeout: 1, // 设置1秒超时
		}

		_, err := rc.Post(slowServer.URL, data, options)
		if err == nil {
			t.Error("预期超时错误但未发生")
		}
	})

	// 测试用例5: 错误URL测试
	t.Run("错误URL测试", func(t *testing.T) {
		data := map[string]interface{}{
			"test": "invalid url",
		}
		options := &PostOptions{}

		_, err := rc.Post("http://invalid.url.test", data, options)
		if err == nil {
			t.Error("预期URL错误但未发生")
		}
	})

	// 测试用例6: 空数据测试
	t.Run("空数据测试", func(t *testing.T) {
		data := map[string]interface{}{}
		options := &PostOptions{}

		resp, err := rc.Post(ts.URL, data, options)
		if err != nil {
			t.Error("POST请求失败:", err)
			return
		}

		if !strings.Contains(resp, "form received") {
			t.Error("响应不符合预期:", resp)
		}
	})
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
