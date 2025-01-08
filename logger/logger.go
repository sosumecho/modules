package logger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogTypeDaily = "daily"
)

var (
	L    *Logger
	once sync.Once
)

type Conf struct {
	Filename  string `mapstructure:"filename"`
	Level     string `mapstructure:"level"`
	Type      string `mapstructure:"type"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxBackup int    `mapstructure:"max_backup"`
	MaxAge    int    `mapstructure:"max_age"`
	Compress  bool   `mapstructure:"compress"`
	IsLocal   bool   `mapstructure:"is_local"`
}

type Logger struct {
	*zap.Logger
	conf *Conf
}

func (l *Logger) Levels() []zapcore.Level {
	return []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}
}

func (l *Logger) getLogWriter(level zapcore.Level) zapcore.WriteSyncer {
	var filename = fmt.Sprintf("%s%s.log", l.conf.Filename, level)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    l.conf.MaxSize,
		MaxBackups: l.conf.MaxBackup,
		MaxAge:     l.conf.MaxAge,
		Compress:   l.conf.Compress,
	}

	if l.conf.IsLocal {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

func (l *Logger) getEncoder() zapcore.Encoder {
	encoderConf := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "Logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     l.customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if l.conf.IsLocal {
		encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConf)
	}
	return zapcore.NewJSONEncoder(encoderConf)
}

func (l *Logger) customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *Logger) setConf(conf *Conf) *Logger {
	l.conf = conf
	return l
}

func (l *Logger) build() *Logger {
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(l.conf.Level)); err != nil {
		panic("日志级别不正确")
	}

	cores := make([]zapcore.Core, 0, len(l.Levels()))
	for _, item := range l.Levels() {
		if item < *level {
			continue
		}
		cores = append(cores, zapcore.NewCore(l.getEncoder(), l.getLogWriter(item), zap.LevelEnablerFunc(func(item zapcore.Level) func(level zapcore.Level) bool {
			return func(level zapcore.Level) bool {
				return level == item
			}
		}(item))))

		//cores = append(cores, zapcore.NewCore(l.getEncoder(), l.getLogWriter(item), item))
	}
	l.Logger = zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		//zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.PanicLevel),
	)
	zap.ReplaceGlobals(l.Logger)
	return l
}

func New(conf *Conf) *Logger {
	return new(Logger).setConf(conf).build()
}
