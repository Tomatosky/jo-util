package fileUtil

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeFixture(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write fixture %q: %v", path, err)
	}
	return path
}

func TestDel(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "nested")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create nested directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("test"), 0o600); err != nil {
		t.Fatalf("create nested file: %v", err)
	}

	if err := Del(dir); err != nil {
		t.Fatalf("Del(%q): %v", dir, err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("Stat(%q) error = %v, want os.ErrNotExist", dir, err)
	}
	if err := Del(dir); err != nil {
		t.Errorf("Del() should be idempotent for a missing path: %v", err)
	}
}

func TestExist(t *testing.T) {
	file := writeFixture(t, "exists.txt", []byte("test"))
	if !Exist(file) {
		t.Errorf("Exist(%q) = false, want true", file)
	}
	if Exist(filepath.Join(t.TempDir(), "missing")) {
		t.Error("Exist() = true for a missing path")
	}
}

func TestGetTotalLines(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{name: "empty", content: "", want: 0},
		{name: "one line without newline", content: "line1", want: 1},
		{name: "three lines", content: "line1\nline2\nline3", want: 3},
		{name: "trailing newline", content: "line1\nline2\n", want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeFixture(t, "lines.txt", []byte(tt.content))
			got, err := GetTotalLines(path)
			if err != nil {
				t.Fatalf("GetTotalLines(%q): %v", path, err)
			}
			if got != tt.want {
				t.Errorf("GetTotalLines(%q) = %d, want %d", tt.content, got, tt.want)
			}
		})
	}

	if got, err := GetTotalLines(filepath.Join(t.TempDir(), "missing")); err == nil || got != 0 {
		t.Errorf("GetTotalLines(missing) = (%d, %v), want (0, error)", got, err)
	}
}

func TestPathTypes(t *testing.T) {
	file := writeFixture(t, "file.txt", []byte("test"))
	dir := t.TempDir()
	missing := filepath.Join(t.TempDir(), "missing")

	tests := []struct {
		name          string
		path          string
		wantDirectory bool
		wantFile      bool
	}{
		{name: "file", path: file, wantFile: true},
		{name: "directory", path: dir, wantDirectory: true},
		{name: "missing", path: missing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDirectory(tt.path); got != tt.wantDirectory {
				t.Errorf("IsDirectory(%q) = %v, want %v", tt.path, got, tt.wantDirectory)
			}
			if got := IsFile(tt.path); got != tt.wantFile {
				t.Errorf("IsFile(%q) = %v, want %v", tt.path, got, tt.wantFile)
			}
		})
	}
}

func TestReadFunctions(t *testing.T) {
	content := []byte("第一行\nsecond line\n")
	path := writeFixture(t, "utf8.txt", content)

	gotBytes, err := ReadBytes(path)
	if err != nil {
		t.Fatalf("ReadBytes(%q): %v", path, err)
	}
	if !bytes.Equal(gotBytes, content) {
		t.Errorf("ReadBytes() = %q, want %q", gotBytes, content)
	}

	gotLines, err := ReadLines(path)
	if err != nil {
		t.Fatalf("ReadLines(%q): %v", path, err)
	}
	if want := []string{"第一行", "second line"}; !reflect.DeepEqual(gotLines, want) {
		t.Errorf("ReadLines() = %#v, want %#v", gotLines, want)
	}

	gotString, err := ReadUtf8String(path)
	if err != nil {
		t.Fatalf("ReadUtf8String(%q): %v", path, err)
	}
	if gotString != string(content) {
		t.Errorf("ReadUtf8String() = %q, want %q", gotString, content)
	}

	missing := filepath.Join(t.TempDir(), "missing")
	if got, err := ReadBytes(missing); err == nil || got != nil {
		t.Errorf("ReadBytes(missing) = (%v, %v), want (nil, error)", got, err)
	}
	if got, err := ReadLines(missing); err == nil || got != nil {
		t.Errorf("ReadLines(missing) = (%v, %v), want (nil, error)", got, err)
	}
	if got, err := ReadUtf8String(missing); err == nil || got != "" {
		t.Errorf("ReadUtf8String(missing) = (%q, %v), want (empty, error)", got, err)
	}
}

func TestRename(t *testing.T) {
	root := t.TempDir()
	oldPath := filepath.Join(root, "old.txt")
	newPath := filepath.Join(root, "new.txt")
	content := []byte("preserved content")
	if err := os.WriteFile(oldPath, content, 0o600); err != nil {
		t.Fatalf("create source file: %v", err)
	}

	if err := Rename(oldPath, newPath); err != nil {
		t.Fatalf("Rename(%q, %q): %v", oldPath, newPath, err)
	}
	if Exist(oldPath) {
		t.Errorf("source path %q still exists", oldPath)
	}
	got, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("read renamed file: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("renamed content = %q, want %q", got, content)
	}
	if err := Rename(oldPath, filepath.Join(root, "other.txt")); err == nil {
		t.Error("Rename() missing source returned nil error")
	}
}

func TestSize(t *testing.T) {
	content := []byte("test content")
	path := writeFixture(t, "size.txt", content)
	size, err := Size(path)
	if err != nil {
		t.Fatalf("Size(%q): %v", path, err)
	}
	if size != int64(len(content)) {
		t.Errorf("Size() = %d, want %d", size, len(content))
	}
	if got, err := Size(filepath.Join(t.TempDir(), "missing")); err == nil || got != 0 {
		t.Errorf("Size(missing) = (%d, %v), want (0, error)", got, err)
	}
}

func TestWriteFunctions(t *testing.T) {
	tests := []struct {
		name  string
		write func(string) error
		want  []byte
	}{
		{name: "bytes", write: func(path string) error { return WriteBytes(path, []byte{0, 1, 2, 255}) }, want: []byte{0, 1, 2, 255}},
		{name: "string", write: func(path string) error { return WriteString(path, "测试 content") }, want: []byte("测试 content")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "output")
			if err := tt.write(path); err != nil {
				t.Fatalf("write %q: %v", path, err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read output: %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("output = %v, want %v", got, tt.want)
			}
		})
	}

	directoryPath := t.TempDir()
	if err := WriteBytes(directoryPath, []byte("x")); err == nil {
		t.Error("WriteBytes(directory) returned nil error")
	}
	if err := WriteString(directoryPath, "x"); err == nil {
		t.Error("WriteString(directory) returned nil error")
	}
}
