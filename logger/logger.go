package logger

import (
	"path"
	"sync"
	"time"

	"github.com/sosumecho/modules/utils"
)

var (
	L    Logger
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

type FieldKind int

const (
	KindAny FieldKind = iota
	KindInt64
	KindString
	KindBool
	KindFloat64
	KindDuration
	KindTime
)

type Field struct {
	Key   string
	Kind  FieldKind
	Int   int64
	Str   string
	Float float64
	Any   interface{}
}

func Int(key string, val int) Field {
	return Field{Key: key, Kind: KindInt64, Int: int64(val)}
}

func Int64(key string, val int64) Field {
	return Field{Key: key, Kind: KindInt64, Int: val}
}

func String(key string, val string) Field {
	return Field{Key: key, Kind: KindString, Str: val}
}

func Bool(key string, val bool) Field {
	var i int64
	if val {
		i = 1
	}
	return Field{Key: key, Kind: KindBool, Int: i}
}

func Float64(key string, val float64) Field {
	return Field{Key: key, Kind: KindFloat64, Float: val}
}

func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Kind: KindDuration, Int: int64(val)}
}

func Time(key string, val time.Time) Field {
	return Field{Key: key, Kind: KindTime, Int: val.UnixNano(), Any: val}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Kind: KindAny, Any: value}
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

func New(conf *Conf) Logger {
	if err := utils.CreateDirIfNotExist(path.Dir(conf.Filename)); err != nil {
		panic(err)
	}
	var l Logger
	switch conf.Type {
	case "zerolog":
		l = newZerologLogger(conf)
	default:
		l = newZapLogger(conf)
	}
	L = l
	return l
}

func SetLogger(l Logger) {
	L = l
}

// Global functions for convenience
func Debug(msg string, fields ...Field) {
	if L != nil {
		L.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...Field) {
	if L != nil {
		L.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...Field) {
	if L != nil {
		L.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...Field) {
	if L != nil {
		L.Error(msg, fields...)
	}
}

func Panic(msg string, fields ...Field) {
	if L != nil {
		L.Panic(msg, fields...)
	}
}

func Fatal(msg string, fields ...Field) {
	if L != nil {
		L.Fatal(msg, fields...)
	}
}
