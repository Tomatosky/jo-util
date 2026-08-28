package httpUtil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type failingReadCloser struct {
	err error
}

func (r failingReadCloser) Read([]byte) (int, error) { return 0, r.err }
func (r failingReadCloser) Close() error             { return nil }

func TestRequestClientConfiguration(t *testing.T) {
	rc := NewRequestClient()
	if rc == nil || rc.Client == nil || rc.Headers == nil {
		t.Fatalf("NewRequestClient() returned incomplete client: %#v", rc)
	}
	transport, ok := rc.Client.Transport.(*http.Transport)
	if !ok || transport.TLSClientConfig == nil || !transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatalf("default transport TLS configuration = %#v", rc.Client.Transport)
	}

	rc.SetTimeout(250 * time.Millisecond)
	if rc.Client.Timeout != 250*time.Millisecond {
		t.Errorf("SetTimeout() = %v, want 250ms", rc.Client.Timeout)
	}

	rc.SetProxy("http://127.0.0.1:8080")
	transport, ok = rc.Client.Transport.(*http.Transport)
	if !ok || transport.Proxy == nil {
		t.Fatalf("SetProxy() transport = %#v, want proxy function", rc.Client.Transport)
	}
	requestURL, _ := url.Parse("https://example.test/resource")
	proxyURL, err := transport.Proxy(&http.Request{URL: requestURL})
	if err != nil {
		t.Fatalf("proxy function: %v", err)
	}
	if got, want := proxyURL.String(), "http://127.0.0.1:8080"; got != want {
		t.Errorf("proxy URL = %q, want %q", got, want)
	}

	defer func() {
		if got := recover(); got == nil {
			t.Error("SetProxy(invalid URL) did not panic")
		}
	}()
	rc.SetProxy("://invalid")
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/resource" || r.URL.Query().Get("q") != "hello world" {
			t.Errorf("request URL = %q, want /resource?q=hello+world", r.URL.String())
		}
		if got := r.Header.Get("X-Test-Header"); got != "test-value" {
			t.Errorf("X-Test-Header = %q, want test-value", got)
		}
		w.Header().Set("X-Response", "response-value")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "response body")
	}))
	t.Cleanup(ts.Close)

	rc := NewRequestClient()
	rc.Headers["X-Test-Header"] = "test-value"
	resp := rc.Get(ts.URL + "/resource?q=hello+world")
	if resp.Err != nil {
		t.Fatalf("Get(): %v", resp.Err)
	}
	if resp.Text() != "response body" {
		t.Errorf("Text() = %q, want response body", resp.Text())
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
	if got := resp.Headers["X-Response"]; got != "response-value" {
		t.Errorf("X-Response = %q, want response-value", got)
	}
}

func TestGetErrors(t *testing.T) {
	if resp := NewRequestClient().Get("://invalid"); resp.Err == nil || len(resp.Body) != 0 {
		t.Errorf("Get(invalid URL) = %#v, want construction error and empty body", resp)
	}

	sentinel := errors.New("transport failed")
	rc := NewRequestClient()
	rc.Client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, sentinel
	})
	if resp := rc.Get("http://example.test"); !errors.Is(resp.Err, sentinel) {
		t.Errorf("Get(transport failure) error = %v, want %v", resp.Err, sentinel)
	}

	sentinel = errors.New("read failed")
	rc.Client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       failingReadCloser{err: sentinel},
		}, nil
	})
	if resp := rc.Get("http://example.test"); !errors.Is(resp.Err, sentinel) {
		t.Errorf("Get(response read failure) error = %v, want %v", resp.Err, sentinel)
	}
}

func TestGetTimeout(t *testing.T) {
	cancelled := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		close(cancelled)
	}))
	t.Cleanup(ts.Close)

	rc := NewRequestClient()
	rc.SetTimeout(50 * time.Millisecond)
	resp := rc.Get(ts.URL)
	if resp.Err == nil {
		t.Fatal("Get() timeout returned nil error")
	}
	var netErr net.Error
	if !errors.As(resp.Err, &netErr) || !netErr.Timeout() {
		t.Errorf("Get() error = %v, want net.Error with Timeout() = true", resp.Err)
	}
	select {
	case <-cancelled:
	case <-time.After(time.Second):
		t.Fatal("server request context was not cancelled after client timeout")
	}
}

func TestRespMethods(t *testing.T) {
	resp := &Resp{Body: []byte(`{"key":"value","num":123}`)}
	if got := resp.Text(); got != `{"key":"value","num":123}` {
		t.Errorf("Text() = %q", got)
	}
	gotMap, err := resp.Json()
	if err != nil {
		t.Fatalf("Json(): %v", err)
	}
	if want := map[string]any{"key": "value", "num": float64(123)}; !reflect.DeepEqual(gotMap, want) {
		t.Errorf("Json() = %#v, want %#v", gotMap, want)
	}

	type payload struct {
		Key string `json:"key"`
		Num int    `json:"num"`
	}
	var gotObject payload
	if err := resp.JsonObj(&gotObject); err != nil {
		t.Fatalf("JsonObj(): %v", err)
	}
	if want := (payload{Key: "value", Num: 123}); gotObject != want {
		t.Errorf("JsonObj() = %#v, want %#v", gotObject, want)
	}

	invalid := &Resp{Body: []byte("not-json")}
	if got, err := invalid.Json(); err == nil || len(got) != 0 {
		t.Errorf("Json(invalid) = (%#v, %v), want empty map and error", got, err)
	}
	if err := invalid.JsonObj(&gotObject); err == nil {
		t.Error("JsonObj(invalid) returned nil error")
	}
}

func TestPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		if got := r.Header.Get("X-Custom"); got != "header-value" {
			t.Errorf("X-Custom = %q, want header-value", got)
		}
		var got map[string]any
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode JSON request: %v", err)
		}
		want := map[string]any{"key": "value", "number": float64(2)}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("JSON body = %#v, want %#v", got, want)
		}
		w.Header().Set("X-Response", "json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, "json response")
	}))
	t.Cleanup(ts.Close)

	rc := NewRequestClient()
	rc.IsJson = true
	rc.Headers["X-Custom"] = "header-value"
	resp := rc.Post(ts.URL, map[string]interface{}{"key": "value", "number": 2})
	assertPostResponse(t, resp, "json response", http.StatusAccepted, "json")
}

func TestPostForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm(): %v", err)
		}
		if want := (url.Values{"key": {"hello world"}, "number": {"2"}}); !reflect.DeepEqual(r.PostForm, want) {
			t.Errorf("form body = %#v, want %#v", r.PostForm, want)
		}
		_, _ = io.WriteString(w, "form response")
	}))
	t.Cleanup(ts.Close)

	resp := NewRequestClient().Post(ts.URL, map[string]interface{}{"key": "hello world", "number": 2})
	if resp.Err != nil || resp.Text() != "form response" || resp.StatusCode != http.StatusOK {
		t.Errorf("Post(form) = %#v", resp)
	}
}

func TestPostMultipart(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "upload name.txt")
	if err := os.WriteFile(filePath, []byte("file content"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
			t.Errorf("Content-Type = %q, want multipart boundary", r.Header.Get("Content-Type"))
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("ParseMultipartForm(): %v", err)
		}
		if got := r.FormValue("field"); got != "value" {
			t.Errorf("field = %q, want value", got)
		}
		upload, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile(): %v", err)
		}
		defer upload.Close()
		if header.Filename != "upload name.txt" {
			t.Errorf("filename = %q, want upload name.txt", header.Filename)
		}
		got, err := io.ReadAll(upload)
		if err != nil {
			t.Fatalf("read upload: %v", err)
		}
		if !bytes.Equal(got, []byte("file content")) {
			t.Errorf("uploaded data = %q", got)
		}
		w.Header().Set("X-Response", "multipart")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "upload response")
	}))
	t.Cleanup(ts.Close)

	rc := NewRequestClient()
	rc.IsMultipart = true
	resp := rc.Post(ts.URL, map[string]interface{}{"file": file, "field": "value"})
	assertPostResponse(t, resp, "upload response", http.StatusCreated, "multipart")
}

func TestPostErrors(t *testing.T) {
	rc := NewRequestClient()
	rc.IsJson = true
	if resp := rc.Post("http://example.test", map[string]interface{}{"invalid": func() {}}); resp.Err == nil {
		t.Error("Post(unencodable JSON) returned nil error")
	}

	rc = NewRequestClient()
	if resp := rc.Post("://invalid", map[string]interface{}{"key": "value"}); resp.Err == nil {
		t.Error("Post(invalid URL) returned nil error")
	}

	sentinel := errors.New("transport failed")
	rc = NewRequestClient()
	rc.Client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, sentinel
	})
	if resp := rc.Post("http://example.test", nil); !errors.Is(resp.Err, sentinel) {
		t.Errorf("Post(transport failure) error = %v, want %v", resp.Err, sentinel)
	}

	sentinel = errors.New("read failed")
	rc.Client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       failingReadCloser{err: sentinel},
		}, nil
	})
	if resp := rc.Post("http://example.test", nil); !errors.Is(resp.Err, sentinel) {
		t.Errorf("Post(response read failure) error = %v, want %v", resp.Err, sentinel)
	}

	closedFile, err := os.CreateTemp(t.TempDir(), "closed")
	if err != nil {
		t.Fatal(err)
	}
	if err := closedFile.Close(); err != nil {
		t.Fatal(err)
	}
	rc = NewRequestClient()
	rc.IsMultipart = true
	if resp := rc.Post("http://example.test", map[string]interface{}{"file": closedFile}); resp.Err == nil {
		t.Error("Post(closed multipart file) returned nil error")
	}
}

func TestPostTimeout(t *testing.T) {
	release := make(chan struct{})
	handlerDone := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(handlerDone)
		<-release
		_, _ = io.WriteString(w, "late response")
	}))
	t.Cleanup(ts.Close)

	rc := NewRequestClient()
	rc.SetTimeout(50 * time.Millisecond)
	resp := rc.Post(ts.URL, map[string]interface{}{"key": "value"})
	if resp.Err == nil {
		t.Fatal("Post() timeout returned nil error")
	}
	var netErr net.Error
	if !errors.As(resp.Err, &netErr) || !netErr.Timeout() {
		t.Errorf("Post() error = %v, want net.Error with Timeout() = true", resp.Err)
	}
	close(release)
	select {
	case <-handlerDone:
	case <-time.After(time.Second):
		t.Fatal("server handler did not exit after timeout test released it")
	}
}

func TestUrlEncode(t *testing.T) {
	params := map[string]interface{}{
		"bool":    true,
		"number":  123,
		"special": "a b&c",
	}
	if got, want := UrlEncode(params), "bool=true&number=123&special=a+b%26c"; got != want {
		t.Errorf("UrlEncode() = %q, want %q", got, want)
	}
	if got := UrlEncode(nil); got != "" {
		t.Errorf("UrlEncode(nil) = %q, want empty string", got)
	}
}

func assertPostResponse(t *testing.T, resp *Resp, body string, status int, header string) {
	t.Helper()
	if resp.Err != nil {
		t.Fatalf("Post(): %v", resp.Err)
	}
	if resp.Text() != body {
		t.Errorf("body = %q, want %q", resp.Text(), body)
	}
	if resp.StatusCode != status {
		t.Errorf("status = %d, want %d", resp.StatusCode, status)
	}
	if got := resp.Headers["X-Response"]; got != header {
		t.Errorf("X-Response = %q, want %q", got, header)
	}
}
