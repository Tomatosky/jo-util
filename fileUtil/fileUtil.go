package fileUtil

import (
	"bufio"
	"os"
)

// Del 删除文件或目录
func Del(path string) error {
	return os.RemoveAll(path)
}

// Exist 判断文件/目录是否存在
func Exist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetTotalLines 获取文件总行数
func GetTotalLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

// IsDirectory 判断是否为目录
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile 判断是否为文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// ReadBytes 读取文件字节内容
func ReadBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadLines 按行读取文件内容
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ReadUtf8String 读取UTF-8格式文件内容
func ReadUtf8String(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Rename 重命名文件/目录
func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Size 获取文件大小
func Size(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// WriteBytes 写入字节内容到文件
func WriteBytes(path string, data []byte) error {
	return os.WriteFile(path, data, 0666)
}

// WriteString 写入字符串到文件
func WriteString(path, content string) error {
	return os.WriteFile(path, []byte(content), 0666)
}
