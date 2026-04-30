package cryptor

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"os"
	"path/filepath"
	"testing"
)

func TestMd5File(t *testing.T) {
	testContent := "hello world"
	expected := computeHash(md5.New(), testContent)

	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(tmpFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	hash, err := Md5File(tmpFile)
	if err != nil {
		t.Errorf("Md5File 失败: %v", err)
	}
	if hash != expected {
		t.Errorf("Md5File 结果错误: expected %s, got %s", expected, hash)
	}

	_, err = Md5File("nonexistent_file.txt")
	if err == nil {
		t.Error("文件不存在时应返回错误")
	}

	tmpDir := t.TempDir()
	hash, err = Md5File(tmpDir)
	if err != nil {
		t.Errorf("Md5File 目录测试失败: %v", err)
	}
	if hash != "" {
		t.Error("传入目录时应返回空字符串")
	}
}

func TestSha1File(t *testing.T) {
	testContent := "hello world"
	expected := computeHash(sha1.New(), testContent)

	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(tmpFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	hash, err := Sha1File(tmpFile)
	if err != nil {
		t.Errorf("Sha1File 失败: %v", err)
	}
	if hash != expected {
		t.Errorf("Sha1File 结果错误: expected %s, got %s", expected, hash)
	}

	_, err = Sha1File("nonexistent_file.txt")
	if err == nil {
		t.Error("文件不存在时应返回错误")
	}

	tmpDir := t.TempDir()
	hash, err = Sha1File(tmpDir)
	if err != nil {
		t.Errorf("Sha1File 目录测试失败: %v", err)
	}
	if hash != "" {
		t.Error("传入目录时应返回空字符串")
	}
}

func TestSha256File(t *testing.T) {
	testContent := "hello world"
	expected := computeHash(sha256.New(), testContent)

	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(tmpFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	hash, err := Sha256File(tmpFile)
	if err != nil {
		t.Errorf("Sha256File 失败: %v", err)
	}
	if hash != expected {
		t.Errorf("Sha256File 结果错误: expected %s, got %s", expected, hash)
	}

	_, err = Sha256File("nonexistent_file.txt")
	if err == nil {
		t.Error("文件不存在时应返回错误")
	}

	tmpDir := t.TempDir()
	hash, err = Sha256File(tmpDir)
	if err != nil {
		t.Errorf("Sha256File 目录测试失败: %v", err)
	}
	if hash != "" {
		t.Error("传入目录时应返回空字符串")
	}
}

func TestSha512File(t *testing.T) {
	testContent := "hello world"
	expected := computeHash(sha512.New(), testContent)

	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(tmpFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	hash, err := Sha512File(tmpFile)
	if err != nil {
		t.Errorf("Sha512File 失败: %v", err)
	}
	if hash != expected {
		t.Errorf("Sha512File 结果错误: expected %s, got %s", expected, hash)
	}

	_, err = Sha512File("nonexistent_file.txt")
	if err == nil {
		t.Error("文件不存在时应返回错误")
	}

	tmpDir := t.TempDir()
	hash, err = Sha512File(tmpDir)
	if err != nil {
		t.Errorf("Sha512File 目录测试失败: %v", err)
	}
	if hash != "" {
		t.Error("传入目录时应返回空字符串")
	}
}

func computeHash(h hash.Hash, content string) string {
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}