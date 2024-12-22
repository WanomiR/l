package l

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	resetColor = "\033[0m"

	LevelDebug  = "DEBUG"
	LevelInfo   = "INFO"
	LevelWarn   = "WARN"
	LevelError  = "ERROR"
	LevelDPanic = "DPANIC"
)

var (
	logger      *zap.Logger
	once        sync.Once
	atomicLevel zap.AtomicLevel
)

// Define ANSI color codes
var levelColors = map[zapcore.Level]string{
	zapcore.DebugLevel:  "\033[36m",   // Cyan
	zapcore.InfoLevel:   "\033[32m",   // Green
	zapcore.WarnLevel:   "\033[33m",   // Yellow
	zapcore.ErrorLevel:  "\033[31m",   // Red
	zapcore.DPanicLevel: "\033[35m",   // Magenta
	zapcore.PanicLevel:  "\033[1;31m", // Bright Red
	zapcore.FatalLevel:  "\033[1;31m", // Bright Red
}

func BuildLogger(logLevel string) {
	once.Do(func() {
		atomicLevel = zap.NewAtomicLevel()
		SetLevel(logLevel)

		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

		// Use the custom level encoder
		cfg.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			color, ok := levelColors[level]
			if !ok {
				color = resetColor
			}
			enc.AppendString(color + level.CapitalString() + resetColor)
		}

		logger = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.AddSync(os.Stdout), atomicLevel), zap.AddCaller())
	})
}

func SetLevel(logLevel string) {
	if logLevel == "" {
		logLevel = "debug"
	}

	switch strings.ToUpper(logLevel) {
	case LevelDebug:
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case LevelInfo:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case LevelWarn:
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case LevelError:
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case LevelDPanic:
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	default:
		panic("invalid log level specified for logger")
	}
}

func CurrentLevel() string {
	return atomicLevel.String()
}

// Logger returns a global logger defined in this package.
// If logger is nil function returns a logger with DEBUG level.
func Logger() *zap.Logger {
	if logger == nil {
		BuildLogger(LevelDebug)
	}

	return logger
}
