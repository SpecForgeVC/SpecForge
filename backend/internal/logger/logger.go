package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var MCP *zap.Logger

func Init() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Backend Logger
	backendFile := &lumberjack.Logger{
		Filename:   "logs/server.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	backendCore := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(backendFile), zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)
	Log = zap.New(backendCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// MCP Logger
	mcpFile := &lumberjack.Logger{
		Filename:   "logs/mcp.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
	}

	mcpCore := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(mcpFile), zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)
	MCP = zap.New(mcpCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

func MCPInfo(msg string, fields ...zap.Field) {
	MCP.Info(msg, fields...)
}

func MCPError(msg string, fields ...zap.Field) {
	MCP.Error(msg, fields...)
}

func MCPWarn(msg string, fields ...zap.Field) {
	MCP.Warn(msg, fields...)
}
