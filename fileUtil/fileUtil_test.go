package fileUtil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDel(t *testing.T) {
	// 创建测试文件
	testFile := "test_del.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}

	// 测试删除文件
	if err := Del(testFile); err != nil {
		t.Error("Del failed:", err)
	}

	// 验证文件是否已删除
	if Exist(testFile) {
		t.Error("File still exists after Del")
	}
}

func TestExist(t *testing.T) {
	// 测试存在的文件
	testFile := "test_exist.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	if !Exist(testFile) {
		t.Error("Exist returned false for existing file")
	}

	// 测试不存在的文件
	if Exist("nonexistent_file.txt") {
		t.Error("Exist returned true for non-existent file")
	}
}

func TestGetTotalLines(t *testing.T) {
	testFile := "test_lines.txt"
	content := "line1\nline2\nline3"
	if err := os.WriteFile(testFile, []byte(content), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	lines, err := GetTotalLines(testFile)
	if err != nil {
		t.Error("GetTotalLines failed:", err)
	}
	if lines != 3 {
		t.Errorf("GetTotalLines returned %d, expected 3", lines)
	}
}

func TestIsDirectory(t *testing.T) {
	// 测试目录
	testDir := "test_dir"
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal("Failed to create test directory:", err)
	}
	defer os.Remove(testDir)

	if !IsDirectory(testDir) {
		t.Error("IsDirectory returned false for directory")
	}

	// 测试文件
	testFile := "test_dir_file.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	if IsDirectory(testFile) {
		t.Error("IsDirectory returned true for file")
	}
}

func TestIsFile(t *testing.T) {
	// 测试文件
	testFile := "test_file.txt"
	if err := os.WriteFile(testFile, []byte("test"), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	if !IsFile(testFile) {
		t.Error("IsFile returned false for file")
	}

	// 测试目录
	testDir := "test_dir"
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal("Failed to create test directory:", err)
	}
	defer os.Remove(testDir)

	if IsFile(testDir) {
		t.Error("IsFile returned true for directory")
	}
}

func TestReadBytes(t *testing.T) {
	testFile := "test_read_bytes.txt"
	content := []byte("test content")
	if err := os.WriteFile(testFile, content, 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	data, err := ReadBytes(testFile)
	if err != nil {
		t.Error("ReadBytes failed:", err)
	}
	if string(data) != string(content) {
		t.Errorf("ReadBytes returned %q, expected %q", data, content)
	}
}

func TestReadLines(t *testing.T) {
	testFile := "test_read_lines.txt"
	content := "line1\nline2\nline3"
	if err := os.WriteFile(testFile, []byte(content), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	lines, err := ReadLines(testFile)
	if err != nil {
		t.Error("ReadLines failed:", err)
	}
	if len(lines) != 3 {
		t.Errorf("ReadLines returned %d lines, expected 3", len(lines))
	}
	if lines[0] != "line1" || lines[1] != "line2" || lines[2] != "line3" {
		t.Errorf("ReadLines returned unexpected content: %v", lines)
	}
}

func TestReadUtf8String(t *testing.T) {
	testFile := "test_read_utf8.txt"
	content := "测试UTF-8内容"
	if err := os.WriteFile(testFile, []byte(content), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	str, err := ReadUtf8String(testFile)
	if err != nil {
		t.Error("ReadUtf8String failed:", err)
	}
	if str != content {
		t.Errorf("ReadUtf8String returned %q, expected %q", str, content)
	}
}

func TestRename(t *testing.T) {
	oldFile := "test_rename_old.txt"
	newFile := "test_rename_new.txt"
	if err := os.WriteFile(oldFile, []byte("test"), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer func() {
		os.Remove(oldFile)
		os.Remove(newFile)
	}()

	if err := Rename(oldFile, newFile); err != nil {
		t.Error("Rename failed:", err)
	}

	if Exist(oldFile) {
		t.Error("Old file still exists after Rename")
	}
	if !Exist(newFile) {
		t.Error("New file does not exist after Rename")
	}
}

func TestSize(t *testing.T) {
	testFile := "test_size.txt"
	content := "test content"
	if err := os.WriteFile(testFile, []byte(content), 0666); err != nil {
		t.Fatal("Failed to create test file:", err)
	}
	defer os.Remove(testFile)

	size, err := Size(testFile)
	if err != nil {
		t.Error("Size failed:", err)
	}
	if size != int64(len(content)) {
		t.Errorf("Size returned %d, expected %d", size, len(content))
	}
}

func TestWriteBytes(t *testing.T) {
	testFile := "test_write_bytes.txt"
	content := []byte("test content")
	defer os.Remove(testFile)

	if err := WriteBytes(testFile, content); err != nil {
		t.Error("WriteBytes failed:", err)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Error("Failed to read test file:", err)
	}
	if string(data) != string(content) {
		t.Errorf("File content %q, expected %q", data, content)
	}
}

func TestWriteString(t *testing.T) {
	testFile := "test_write_string.txt"
	content := "test content"
	defer os.Remove(testFile)

	if err := WriteString(testFile, content); err != nil {
		t.Error("WriteString failed:", err)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Error("Failed to read test file:", err)
	}
	if string(data) != content {
		t.Errorf("File content %q, expected %q", data, content)
	}
}

func TestCleanup(t *testing.T) {
	// 清理可能遗留的测试文件
	files, _ := filepath.Glob("test_*")
	for _, f := range files {
		os.Remove(f)
	}
}
