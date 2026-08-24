package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
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
