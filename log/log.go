package log

import (
	"bufio"
	"fmt"
	"github.com/sencydai/utils"
	"os"
	"path"
	"time"
)

type LogLevel = int

//日志等级
const (
	DEBUG_N LogLevel = iota
	INFO_N
	WARN_N
	ERROR_N
	FATAL_N
	TRASH
)

const (
	sDEBUG = "DEBUG"
	sINFO  = "INFO"
	sWARN  = "WARN"
	sERROR = "ERROR"
	sFATAL = "FATAL"

	syncPeriod     = time.Millisecond * 100
	defaultBufSize = 1024 * 1024
)

var (
	levelText = map[int]string{DEBUG_N: sDEBUG, INFO_N: sINFO, WARN_N: sWARN, ERROR_N: sERROR, FATAL_N: sFATAL}

	loggerMgr = make(map[string]*loggerData)
)

type loggerData struct {
	level  LogLevel
	file   *os.File
	writer *bufio.Writer

	chOutput chan string
}

var DefaultLogger = GetLogger("DefaultLogger")

type ILogger interface {
	SetFileName(fileName string) error
	SetLevel(level LogLevel) bool
	Close()

	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

func newLogger() *loggerData {
	logger := &loggerData{level: DEBUG_N, chOutput: make(chan string, 100)}
	go func() {
		var output string
		var hasData bool
		write := os.Stdout.WriteString
		for {
			select {
			case output = <-logger.chOutput:
				write(output)
				if logger.writer != nil {
					logger.writer.WriteString(output)
					hasData = true
				}
			case <-time.After(syncPeriod):
				if hasData {
					logger.writer.Flush()
					logger.file.Sync()
					hasData = false
				}
			}
		}
	}()

	return logger
}

func GetLogger(name string) ILogger {
	if logger, ok := loggerMgr[name]; ok {
		return logger
	}
	loggerMgr[name] = newLogger()
	return loggerMgr[name]
}

func (l *loggerData) SetFileName(fileName string) error {
	os.MkdirAll(path.Dir(fileName), os.ModeDir)
	if file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		return err
	} else {
		l.file, l.writer = file, bufio.NewWriterSize(file, defaultBufSize)
		return nil
	}
}

func (l *loggerData) SetLevel(level LogLevel) bool {
	if level < DEBUG_N || level >= TRASH {
		return false
	}

	l.level = level
	return true
}

func (l *loggerData) Close() {
	for {
		if len(l.chOutput) == 0 {
			time.Sleep(syncPeriod + time.Millisecond*10)
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (l *loggerData) writeBufferf(level LogLevel, skip int, format string, data ...interface{}) {
	if level >= l.level {
		l.chOutput <- fmt.Sprintf("%s %s [%s] - %s\n", utils.TimeFormat(time.Now()), levelText[level], utils.FileLine(skip), fmt.Sprintf(format, data...))
	}
}

func (l *loggerData) writeBuffer(level LogLevel, skip int, data ...interface{}) {
	if level >= l.level {
		l.chOutput <- fmt.Sprintf("%s %s [%s] - %s\n", utils.TimeFormat(time.Now()), levelText[level], utils.FileLine(skip), fmt.Sprint(data...))
	}
}

func (l *loggerData) Debug(data ...interface{}) {
	l.writeBuffer(DEBUG_N, 3, data...)
}

func (l *loggerData) Debugf(format string, data ...interface{}) {
	l.writeBufferf(DEBUG_N, 3, format, data...)
}

func (l *loggerData) Info(data ...interface{}) {
	l.writeBuffer(INFO_N, 3, data...)
}

func (l *loggerData) Infof(format string, data ...interface{}) {
	l.writeBufferf(INFO_N, 3, format, data...)
}

func (l *loggerData) Warn(data ...interface{}) {
	l.writeBuffer(WARN_N, 3, data...)
}

func (l *loggerData) Warnf(format string, data ...interface{}) {
	l.writeBufferf(WARN_N, 3, format, data...)
}

func (l *loggerData) Error(data ...interface{}) {
	l.writeBuffer(ERROR_N, 3, data...)
}

func (l *loggerData) Errorf(format string, data ...interface{}) {
	l.writeBufferf(ERROR_N, 3, format, data...)
}

func (l *loggerData) Fatal(data ...interface{}) {
	l.writeBuffer(FATAL_N, 3, data...)
}

func (l *loggerData) Fatalf(format string, data ...interface{}) {
	l.writeBufferf(FATAL_N, 3, format, data...)
}
