package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitLogWriterFormatMatchesConsoleWithoutColor(t *testing.T) {
	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = originalStdout
		_ = reader.Close()
		_ = writer.Close()
	})

	var fileOutput bytes.Buffer
	log := InitLog(map[io.Writer]zapcore.Level{
		&fileOutput: zapcore.InfoLevel,
	})
	log.Info("application started", zap.String("module", "main"))

	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	os.Stdout = originalStdout

	consoleOutput, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read console output: %v", err)
	}
	wantConsoleOutput := "\x1b[34m" + fileOutput.String() + "\x1b[0m"
	if string(consoleOutput) != wantConsoleOutput {
		t.Errorf("file output does not match console output without color\nconsole: %q\nfile: %q", consoleOutput, fileOutput.String())
	}
	if strings.Contains(fileOutput.String(), "\x1b[") {
		t.Errorf("file output contains ANSI color codes: %q", fileOutput.String())
	}
	if json.Valid(fileOutput.Bytes()) {
		t.Errorf("file output is still JSON: %s", fileOutput.String())
	}
}

func TestGetColor(t *testing.T) {
	tests := []struct {
		level zapcore.Level
		want  string
	}{
		{level: zapcore.DebugLevel, want: "\x1b[90m"},
		{level: zapcore.InfoLevel, want: "\x1b[34m"},
		{level: zapcore.WarnLevel, want: "\x1b[33m"},
		{level: zapcore.ErrorLevel, want: "\x1b[31m"},
		{level: zapcore.DPanicLevel, want: "\x1b[31m"},
		{level: zapcore.PanicLevel, want: "\x1b[31m"},
		{level: zapcore.FatalLevel, want: "\x1b[31m"},
		{level: zapcore.Level(99), want: "\x1b[0m"},
	}
	for _, tt := range tests {
		if got := getColor(tt.level); got != tt.want {
			t.Errorf("getColor(%v) = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestColorEncoderClonePreservesColor(t *testing.T) {
	config := zap.NewProductionEncoderConfig()
	clone := (&ColorEncoder{Encoder: zapcore.NewConsoleEncoder(config)}).Clone()
	colorClone, ok := clone.(*ColorEncoder)
	if !ok {
		t.Fatalf("ColorEncoder.Clone() returned %T, want *ColorEncoder", clone)
	}
	buf, err := colorClone.EncodeEntry(zapcore.Entry{Level: zapcore.WarnLevel, Message: "warning"}, nil)
	if err != nil {
		t.Fatalf("EncodeEntry(): %v", err)
	}
	defer buf.Free()
	if got := buf.String(); !strings.HasPrefix(got, "\x1b[33m") || !strings.HasSuffix(got, "\x1b[0m") {
		t.Errorf("cloned encoder output = %q, want warning color wrapper", got)
	}
}

func TestInitLogWriterLevelFiltering(t *testing.T) {
	var infoOutput bytes.Buffer
	var errorOutput bytes.Buffer
	log := InitLog(map[io.Writer]zapcore.Level{
		&infoOutput:  zapcore.InfoLevel,
		&errorOutput: zapcore.ErrorLevel,
	})
	log.Debug("debug-message")
	log.Info("info-message")
	log.Warn("warn-message")
	log.Error("error-message")

	infoText := infoOutput.String()
	if strings.Contains(infoText, "debug-message") || !strings.Contains(infoText, "info-message") || !strings.Contains(infoText, "warn-message") || !strings.Contains(infoText, "error-message") {
		t.Errorf("info writer output has wrong level set: %q", infoText)
	}
	errorText := errorOutput.String()
	if strings.Contains(errorText, "debug-message") || strings.Contains(errorText, "info-message") || strings.Contains(errorText, "warn-message") || !strings.Contains(errorText, "error-message") {
		t.Errorf("error writer output has wrong level set: %q", errorText)
	}
}

func TestSimplyInit(t *testing.T) {
	logDir := filepath.Join(t.TempDir(), "nested", "logs")
	log := SimplyInit(logDir)
	if log == nil {
		t.Fatal("SimplyInit() returned nil")
	}
	info, err := os.Stat(logDir)
	if err != nil {
		t.Fatalf("Stat(log directory): %v", err)
	}
	if !info.IsDir() {
		t.Errorf("SimplyInit() path is not a directory: %v", info.Mode())
	}

	filePath := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(filePath, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := loggerPanicValue(func() { SimplyInit(filePath) }); got != "logPath exists but is not a directory: "+filePath {
		t.Errorf("SimplyInit(file) panic = %v", got)
	}
	if got := loggerPanicValue(func() { SimplyInit(filepath.Join(filePath, "logs")) }); got == nil {
		t.Error("SimplyInit(path below a file) did not propagate directory creation error")
	}
}

func loggerPanicValue(f func()) (recovered any) {
	defer func() { recovered = recover() }()
	f()
	return nil
}
