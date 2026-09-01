package tgk

import (
	"io"
	"strings"
	"sync"
)

// LogLevel はログレベルを表す型です
type LogLevel int

const (
	// LogLevelDebug はデバッグレベルです
	LogLevelDebug LogLevel = iota
	// LogLevelInfo は情報レベルです
	LogLevelInfo
	// LogLevelWarn は警告レベルです
	LogLevelWarn
	// LogLevelError はエラーレベルです
	LogLevelError
)

// String はログレベルの文字列表現を返します
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLogLevel は文字列をLogLevelに解析します
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return LogLevelInfo // デフォルトはInfoレベル
	}
}

// LoggerService はログサービスのインターフェースです
type LoggerService interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	SetLevel(level LogLevel)
	SetOutput(writer io.Writer)
	AddOutput(writer io.Writer)
}

// loggerService はログサービスの実装です
type loggerService struct {
	level LogLevel
	mu    sync.Mutex
}

// NewLoggerService はログサービスを作成します
func NewLoggerService() LoggerService {
	return &loggerService{
		level: LogLevelInfo, // デフォルトはInfoレベル
	}
}

// formatMessage はメッセージをフォーマットします
func formatMessage(level LogLevel, format string, args ...interface{}) string {
	prefix := "[" + level.String() + "] "

	// argsがない場合は単純に連結
	if len(args) == 0 {
		return prefix + format
	}

	// 簡易的なフォーマット処理
	result := format
	for _, arg := range args {
		// %vを見つけて置換
		pos := strings.Index(result, "%")
		if pos >= 0 && pos+1 < len(result) {
			switch result[pos+1] {
			case 'd', 'v', 's', 't':
				var value string
				switch v := arg.(type) {
				case string:
					value = v
				case int:
					value = intToString(v)
				case bool:
					if v {
						value = "true"
					} else {
						value = "false"
					}
				default:
					value = "<unknown>"
				}
				result = result[:pos] + value + result[pos+2:]
			}
		}
	}

	return prefix + result
}

// intToString は整数を文字列に変換します
func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	negative := false
	if n < 0 {
		negative = true
		n = -n
	}

	var digits [20]byte // 64ビット整数の最大桁数
	i := len(digits) - 1

	for n > 0 {
		digits[i] = byte('0' + n%10)
		n /= 10
		i--
	}

	if negative {
		digits[i] = '-'
		i--
	}

	return string(digits[i+1:])
}

// log は共通のログ出力メソッドです
func (l *loggerService) log(level LogLevel, format string, args ...interface{}) {
	if l.level <= level {
		message := formatMessage(level, format, args...)

		// println()を使用してログを出力
		println(message)
	}
}

func (l *loggerService) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

func (l *loggerService) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

func (l *loggerService) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

func (l *loggerService) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

func (l *loggerService) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput は出力先を設定します（TinyGoでは使用しないため空実装）
func (l *loggerService) SetOutput(writer io.Writer) {
	// TinyGoではprintln()を使用するため、この関数は何もしない
}

// AddOutput は出力先を追加します（TinyGoでは使用しないため空実装）
func (l *loggerService) AddOutput(writer io.Writer) {
	// TinyGoではprintln()を使用するため、この関数は何もしない
}
