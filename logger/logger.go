package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func init() {
	Log = InitLog(nil)
}

func SimplyInit(logPath string) *zap.Logger {
	info, err := os.Stat(logPath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if err == nil && !info.IsDir() {
		panic("logPath exists but is not a directory: " + logPath)
	}
	if os.IsNotExist(err) {
		_ = os.MkdirAll(logPath, 0755)
	}
	infoPath := logPath + "/info.log"
	warnPath := logPath + "/warn.log"
	errorPath := logPath + "/error.log"
	return InitLog(map[io.Writer]zapcore.Level{
		&lumberjack.Logger{
			Filename:   infoPath, //日志文件存放目录，如果文件夹不存在会自动创建
			MaxSize:    50,       //文件大小限制,单位MB
			MaxBackups: 10,       //最大保留日志文件数量
			MaxAge:     30,       //日志文件保留天数
			LocalTime:  true,
		}: zapcore.InfoLevel,
		&lumberjack.Logger{
			Filename:   warnPath, //日志文件存放目录，如果文件夹不存在会自动创建
			MaxSize:    50,       //文件大小限制,单位MB
			MaxBackups: 10,       //最大保留日志文件数量
			MaxAge:     60,       //日志文件保留天数
			LocalTime:  true,
		}: zapcore.WarnLevel,
		&lumberjack.Logger{
			Filename:   errorPath, //日志文件存放目录，如果文件夹不存在会自动创建
			MaxSize:    50,        //文件大小限制,单位MB
			MaxBackups: 10,        //最大保留日志文件数量
			MaxAge:     90,        //日志文件保留天数
			LocalTime:  true,
		}: zapcore.ErrorLevel,
	})
}

func InitLog(writers map[io.Writer]zapcore.Level) *zap.Logger {
	var coreArr []zapcore.Core

	//日志级别
	debugLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev <= zap.FatalLevel && lev >= zap.DebugLevel
	})

	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") //指定时间格式
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	colorEncoder := &ColorEncoder{Encoder: consoleEncoder}
	consoleCore := zapcore.NewCore(colorEncoder, zapcore.AddSync(os.Stdout), debugLevel)
	coreArr = append(coreArr, consoleCore)

	if writers != nil && len(writers) > 0 {
		// 添加多个 writer
		fileEncoderConfig := zap.NewProductionEncoderConfig()
		fileEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

		for w, level := range writers {
			levelEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
				return lev <= zap.FatalLevel && lev >= level
			})
			fileWriteSyncer := zapcore.AddSync(w)
			fileCore := zapcore.NewCore(fileEncoder, fileWriteSyncer, levelEnabler)
			coreArr = append(coreArr, fileCore)
		}
	}

	return zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}
