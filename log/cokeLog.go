package log

import (
	"github.com/chuccp/utils/queue"
	"log"
	"time"
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger struct {
	config *Config
	queue  *queue.Queue
}

func New() *Logger {
	log := &Logger{config: defaultConfig, queue: queue.NewQueue()}
	go log.init()
	return log
}
func (logger *Logger) init() {
	for {
		v, _ := logger.queue.Poll()
		p := v.(*Entry)
		p.WriteTo()
		freeEntry(p)
	}
}
func (logger *Logger) Info(format string, args ...interface{}) {
	var ti *time.Time
	if logger.config.level >= InfoLevel {
		now := time.Now()
		ti = &now
		p := newEntry(logger.config, logger.config.Out)
		p.Log(InfoLevel, &now, format, args...)
		logger.queue.Offer(p)
	}
	logger.writeFile(InfoLevel, ti, format, args...)
}
func (logger *Logger) Debug(format string, args ...interface{}) {
	var ti *time.Time
	if logger.config.level >= DebugLevel {
		now := time.Now()
		ti = &now
		p := newEntry(logger.config, logger.config.Out)
		p.Log(DebugLevel, &now, format, args...)
		logger.queue.Offer(p)
	}
	logger.writeFile(DebugLevel, ti, format, args...)
}
func (logger *Logger) Trace(format string, args ...interface{}) {
	var ti *time.Time
	if logger.config.level >= TraceLevel {
		now := time.Now()
		ti = &now
		p := newEntry(logger.config, logger.config.Out)
		p.Log(TraceLevel, &now, format, args...)
		logger.queue.Offer(p)
	}
	logger.writeFile(TraceLevel, ti, format, args...)
}
func (logger *Logger) Fatal(format string, args ...interface{}) {
	var ti *time.Time
	if logger.config.level >= FatalLevel {
		now := time.Now()
		ti = &now
		p := newEntry(logger.config, logger.config.Out)
		p.Log(FatalLevel, &now, format, args...)
		logger.queue.Offer(p)
	}
	logger.writeFile(FatalLevel, ti, format, args...)
}
func (logger *Logger) Panic(format string, args ...interface{}) {
	var ti *time.Time
	if logger.config.level >= PanicLevel {
		now := time.Now()
		ti = &now
		p := newEntry(logger.config, logger.config.Out)
		p.Log(PanicLevel, &now, format, args...)
		logger.queue.Offer(p)
	}
	logger.writeFile(PanicLevel, ti, format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	var ti *time.Time
	if logger.config.level >= ErrorLevel {
		now := time.Now()
		ti = &now
		p := newEntry(logger.config, logger.config.Out)
		p.Log(ErrorLevel, &now, format, args...)
		logger.queue.Offer(p)
	}
	logger.writeFile(ErrorLevel, ti, format, args...)
}

func (logger *Logger) writeFile(level Level, ti *time.Time, format string, args ...interface{}) {
	ct := logger.config.levelMap[level]
	if ct != nil {
		if ti == nil {
			now := time.Now()
			ti = &now
		}
		w, err := ct.getOut(ti)
		if err == nil {
			p := newEntry(logger.config, w)
			p.Log(level, ti, format, args...)
			logger.queue.Offer(p)
		} else {
			log.Println(err)
		}
	}
}

var defaultLogger = New()

func InfoF(format string, value ...interface{}) {
	defaultLogger.Info(format, value...)
}
func DebugF(format string, value ...interface{}) {
	defaultLogger.Debug(format, value...)
}
func FatalF(format string, value ...interface{}) {
	defaultLogger.Fatal(format, value...)
}
func ErrorF(format string, value ...interface{}) {
	defaultLogger.Error(format, value...)
}
func TraceF(format string, value ...interface{}) {
	defaultLogger.Trace(format, value...)
}
func PanicF(format string, value ...interface{}) {
	defaultLogger.Panic(format, value...)
}

func Info(value ...interface{}) {
	defaultLogger.Info("", value...)
}
func Debug(value ...interface{}) {
	defaultLogger.Debug("", value...)
}
func Fatal(value ...interface{}) {
	defaultLogger.Fatal("", value...)
}
func Error(value ...interface{}) {
	defaultLogger.Error("", value...)
}
func Trace(value ...interface{}) {
	defaultLogger.Trace("", value...)
}
func Panic(value ...interface{}) {
	defaultLogger.Panic("", value...)
}
