package logger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	*zap.Logger
	conf *Conf
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, l.toZapFields(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, l.toZapFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, l.toZapFields(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, l.toZapFields(fields)...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.Logger.Panic(msg, l.toZapFields(fields)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.Logger.Fatal(msg, l.toZapFields(fields)...)
}

func (l *zapLogger) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch f.Kind {
		case KindInt64:
			zapFields[i] = zap.Int64(f.Key, f.Int)
		case KindString:
			zapFields[i] = zap.String(f.Key, f.Str)
		case KindBool:
			zapFields[i] = zap.Bool(f.Key, f.Int == 1)
		case KindFloat64:
			zapFields[i] = zap.Float64(f.Key, f.Float)
		case KindDuration:
			zapFields[i] = zap.Duration(f.Key, time.Duration(f.Int))
		case KindTime:
			if t, ok := f.Any.(time.Time); ok {
				zapFields[i] = zap.Time(f.Key, t)
			} else {
				zapFields[i] = zap.Time(f.Key, time.Unix(0, f.Int))
			}
		default:
			zapFields[i] = zap.Any(f.Key, f.Any)
		}
	}
	return zapFields
}

func (l *zapLogger) levels() []zapcore.Level {
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

func (l *zapLogger) getLogWriter(level zapcore.Level) zapcore.WriteSyncer {
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

func (l *zapLogger) getEncoder() zapcore.Encoder {
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

func (l *zapLogger) customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *zapLogger) build() *zapLogger {
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(l.conf.Level)); err != nil {
		panic("日志级别不正确")
	}

	cores := make([]zapcore.Core, 0, len(l.levels()))
	for _, item := range l.levels() {
		if item < *level {
			continue
		}
		cores = append(cores, zapcore.NewCore(l.getEncoder(), l.getLogWriter(item), zap.LevelEnablerFunc(func(item zapcore.Level) func(level zapcore.Level) bool {
			return func(level zapcore.Level) bool {
				return level == item
			}
		}(item))))
	}
	l.Logger = zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zap.PanicLevel),
	)
	zap.ReplaceGlobals(l.Logger)
	return l
}

func newZapLogger(conf *Conf) Logger {
	l := &zapLogger{conf: conf}
	return l.build()
}
