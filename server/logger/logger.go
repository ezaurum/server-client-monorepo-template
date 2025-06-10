package logger

// logger 패키지
// 이 패키지는 기본 로거를 설정하고, 다양한 로그 레벨을 지원하는 함수를 제공합니다.
// 로그는 JSON 형식으로 출력됩니다.
// logger 패키지를 사용하여 로그 메시지를 기록할 수 있습니다.

import (
	"log/slog"
	"os"
)

var _logger *slog.Logger

func init() {
	// 기본 로거 설정
	_logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// Logger Get default logger
func Logger() *slog.Logger {
	return _logger
}

func Debug(msg string, v ...any) {
	_logger.Debug(msg, v...)
}

func Info(msg string, v ...any) {
	_logger.Info(msg, v...)
}

func Warn(msg string, v ...any) {
	_logger.Warn(msg, v...)
}

func Error(msg string, v ...any) {
	_logger.Error(msg, v...)
}
