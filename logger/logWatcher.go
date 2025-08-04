package logger

//
//import (
//	"bufio"
//	"fmt"
//	"github.com/fsnotify/fsnotify"
//	"io"
//	"os"
//	"path/filepath"
//	"regexp"
//	"strings"
//)
//
//type LogWatcher struct {
//	watcher      *fsnotify.Watcher // fsnotify 监视器
//	files        map[string]*logWatcher
//	directory    string
//	patternRegex *regexp.Regexp
//}
//
//type logWatcher struct {
//	file    *os.File
//	lastPos int64
//}
//
//func NewLogWatcher(pathPattern string) *LogWatcher {
//	if Log == nil {
//		panic("init log first")
//	}
//
//	// 分离目录和文件模式
//	directory, pattern := filepath.Split(pathPattern)
//	if directory == "" {
//		directory = "."
//	}
//
//	// 转换通配符模式为正则表达式
//	regexPattern := "^" + strings.ReplaceAll(
//		strings.ReplaceAll(pattern, ".", `\.`),
//		"*", ".*") + "$"
//
//	patternRegex, err := regexp.Compile(regexPattern)
//	if err != nil {
//		panic(fmt.Errorf("无效的文件模式: %v", err))
//	}
//
//	// 创建 fsnotify 监视器
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		panic(fmt.Errorf("无法创建文件监视器: %v", err))
//	}
//
//	// 添加目录监视
//	err = watcher.Add(directory)
//	if err != nil {
//		_ = watcher.Close()
//		panic(fmt.Errorf("无法监视目录: %v", err))
//	}
//
//	lw := &LogWatcher{
//		watcher:      watcher,
//		files:        make(map[string]*logWatcher),
//		directory:    directory,
//		patternRegex: patternRegex,
//	}
//
//	// 初始扫描目录
//	err = lw.scanDirectory()
//	if err != nil {
//		_ = watcher.Close()
//		panic(err)
//	}
//
//	return lw
//}
//
//func (w *LogWatcher) scanDirectory() error {
//	// 读取目录内容
//	entries, err := os.ReadDir(w.directory)
//	if err != nil {
//		return fmt.Errorf("无法读取目录: %v", err)
//	}
//
//	// 检查目录中的每个文件
//	for _, entry := range entries {
//		if entry.IsDir() {
//			continue
//		}
//
//		filename := entry.Name()
//		if !w.patternRegex.MatchString(filename) {
//			continue
//		}
//
//		fullPath := filepath.Join(w.directory, filename)
//		if _, exists := w.files[fullPath]; exists {
//			continue
//		}
//
//		// 新文件，添加到监视列表
//		fw, err := newLogWatcher(fullPath)
//		if err != nil {
//			Log.Warn(fmt.Sprintf("无法监视文件 %s: %v", fullPath, err.Error()))
//			continue
//		}
//		w.files[fullPath] = fw
//	}
//
//	return nil
//}
//
//func newLogWatcher(filename string) (*logWatcher, error) {
//	file, err := os.Open(filename)
//	if err != nil {
//		return nil, fmt.Errorf("无法打开文件: %v", err)
//	}
//
//	pos, err := file.Seek(0, io.SeekEnd)
//	if err != nil {
//		_ = file.Close()
//		return nil, fmt.Errorf("无法定位到文件末尾: %v", err)
//	}
//
//	return &logWatcher{
//		file:    file,
//		lastPos: pos,
//	}, nil
//}
//
//func (w *LogWatcher) Watch() (map[string][]string, error) {
//	results := make(map[string][]string) //map[filename]lines
//
//	select {
//	case event, ok := <-w.watcher.Events:
//		if !ok {
//			return nil, fmt.Errorf("监视器通道已关闭")
//		}
//		Log.Info(fmt.Sprintf("file watcher event: %v", event))
//
//		// 处理不同事件类型
//		switch {
//		case event.Op&fsnotify.Create == fsnotify.Create:
//			// 新文件创建
//			if _, exists := w.files[event.Name]; !exists && w.patternRegex.MatchString(filepath.Base(event.Name)) {
//				_, err := newLogWatcher(event.Name)
//				if err != nil {
//					Log.Warn(fmt.Sprintf("无法监视新文件 %s: %v", event.Name, err.Error()))
//				}
//			}
//
//		case event.Op&fsnotify.Write == fsnotify.Write:
//			// 文件写入
//			if fw, exists := w.files[event.Name]; exists {
//				lines, err := w.readNewLines(fw)
//				if err != nil {
//					return nil, fmt.Errorf("读取文件 %s 新内容失败: %v", event.Name, err)
//				}
//				if len(lines) > 0 {
//					results[event.Name] = lines
//				}
//			}
//
//		case event.Op&fsnotify.Remove == fsnotify.Remove:
//			// 文件删除
//			if fw, exists := w.files[event.Name]; exists {
//				_ = fw.file.Close()
//				delete(w.files, event.Name)
//			}
//
//		case event.Op&fsnotify.Rename == fsnotify.Rename:
//			// 文件重命名
//			if fw, exists := w.files[event.Name]; exists {
//				_ = fw.file.Close()
//				delete(w.files, event.Name)
//			}
//		}
//
//	case err, ok := <-w.watcher.Errors:
//		if !ok {
//			return nil, fmt.Errorf("监视器错误通道已关闭")
//		}
//		return nil, fmt.Errorf("监视器错误: %v", err)
//	default:
//		return results, nil
//	}
//
//	return results, nil
//}
//
//func (w *LogWatcher) readAllLines(fw *logWatcher) ([]string, error) {
//	var lines []string
//
//	_, err := fw.file.Seek(0, io.SeekStart)
//	if err != nil {
//		return nil, fmt.Errorf("无法定位到文件开头: %v", err)
//	}
//
//	scanner := bufio.NewScanner(fw.file)
//	for scanner.Scan() {
//		lines = append(lines, scanner.Text())
//	}
//
//	if err = scanner.Err(); err != nil {
//		return lines, fmt.Errorf("读取文件错误: %v", err)
//	}
//
//	// 更新最后位置
//	fw.lastPos, err = fw.file.Seek(0, io.SeekCurrent)
//	if err != nil {
//		return lines, fmt.Errorf("无法获取当前位置: %v", err)
//	}
//
//	return lines, nil
//}
//
//func (w *LogWatcher) readNewLines(fw *logWatcher) ([]string, error) {
//	var lines []string
//
//	info, err := fw.file.Stat()
//	if err != nil {
//		return nil, fmt.Errorf("无法获取文件状态: %v", err)
//	}
//
//	currentSize := info.Size()
//
//	if currentSize > fw.lastPos {
//		// 文件有新增内容
//		_, err := fw.file.Seek(fw.lastPos, io.SeekStart)
//		if err != nil {
//			return nil, fmt.Errorf("无法定位到上次位置: %v", err)
//		}
//
//		scanner := bufio.NewScanner(fw.file)
//		for scanner.Scan() {
//			lines = append(lines, scanner.Text())
//		}
//
//		if err = scanner.Err(); err != nil {
//			return lines, fmt.Errorf("读取文件错误: %v", err)
//		}
//
//		// 更新最后位置
//		fw.lastPos, err = fw.file.Seek(0, io.SeekCurrent)
//		if err != nil {
//			return lines, fmt.Errorf("无法获取当前位置: %v", err)
//		}
//	} else if currentSize < fw.lastPos {
//		// 文件被截断，重新读取整个文件
//		_, err := fw.file.Seek(0, io.SeekStart)
//		if err != nil {
//			return nil, fmt.Errorf("无法定位到文件开头: %v", err)
//		}
//
//		scanner := bufio.NewScanner(fw.file)
//		for scanner.Scan() {
//			lines = append(lines, scanner.Text())
//		}
//
//		if err = scanner.Err(); err != nil {
//			return lines, fmt.Errorf("读取文件错误: %v", err)
//		}
//
//		// 更新最后位置
//		fw.lastPos, err = fw.file.Seek(0, io.SeekCurrent)
//		if err != nil {
//			return lines, fmt.Errorf("无法获取当前位置: %v", err)
//		}
//	}
//
//	return lines, nil
//}
//
//func (w *LogWatcher) Close() error {
//	var errs []error
//
//	// 关闭所有文件
//	for _, fw := range w.files {
//		if err := fw.file.Close(); err != nil {
//			errs = append(errs, err)
//		}
//	}
//
//	// 关闭 fsnotify 监视器
//	if err := w.watcher.Close(); err != nil {
//		errs = append(errs, err)
//	}
//
//	if len(errs) > 0 {
//		return fmt.Errorf("关闭时出错: %v", errs)
//	}
//	return nil
//}
