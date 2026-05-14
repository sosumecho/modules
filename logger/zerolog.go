package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zerologLogger struct {
	logger zerolog.Logger
	conf   *Conf
}

func (l *zerologLogger) Debug(msg string, fields ...Field) {
	l.appendFields(l.logger.Debug(), fields).Msg(msg)
}

func (l *zerologLogger) Info(msg string, fields ...Field) {
	l.appendFields(l.logger.Info(), fields).Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...Field) {
	l.appendFields(l.logger.Warn(), fields).Msg(msg)
}

func (l *zerologLogger) Error(msg string, fields ...Field) {
	l.appendFields(l.logger.Error(), fields).Msg(msg)
}

func (l *zerologLogger) Panic(msg string, fields ...Field) {
	l.appendFields(l.logger.Panic(), fields).Msg(msg)
}

func (l *zerologLogger) Fatal(msg string, fields ...Field) {
	l.appendFields(l.logger.Fatal(), fields).Msg(msg)
}

func (l *zerologLogger) appendFields(event *zerolog.Event, fields []Field) *zerolog.Event {
	for _, f := range fields {
		switch f.Kind {
		case KindInt64:
			event = event.Int64(f.Key, f.Int)
		case KindString:
			event = event.Str(f.Key, f.Str)
		case KindBool:
			event = event.Bool(f.Key, f.Int == 1)
		case KindFloat64:
			event = event.Float64(f.Key, f.Float)
		case KindDuration:
			event = event.Dur(f.Key, time.Duration(f.Int))
		case KindTime:
			if t, ok := f.Any.(time.Time); ok {
				event = event.Time(f.Key, t)
			} else {
				event = event.Time(f.Key, time.Unix(0, f.Int))
			}
		default:
			event = event.Interface(f.Key, f.Any)
		}
	}
	return event
}

func (l *zerologLogger) build() *zerologLogger {
	level, err := zerolog.ParseLevel(l.conf.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	var writer io.Writer
	if l.conf.Filename != "" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   l.conf.Filename + ".log",
			MaxSize:    l.conf.MaxSize,
			MaxBackups: l.conf.MaxBackup,
			MaxAge:     l.conf.MaxAge,
			Compress:   l.conf.Compress,
		}
		if l.conf.IsLocal {
			writer = zerolog.MultiLevelWriter(os.Stdout, lumberJackLogger)
		} else {
			writer = lumberJackLogger
		}
	} else {
		writer = os.Stdout
	}

	if l.conf.IsLocal {
		writer = zerolog.ConsoleWriter{Out: writer, TimeFormat: "2006-01-02 15:04:05"}
	}

	l.logger = zerolog.New(writer).Level(level).With().Timestamp().Caller().Logger()
	log.Logger = l.logger
	return l
}

func newZerologLogger(conf *Conf) Logger {
	l := &zerologLogger{conf: conf}
	return l.build()
}
