package log

import (
	"github.com/go-xorm/core"
)

type SqlLogger struct {
	logger  *loggerData
	showSql bool
	level   core.LogLevel
}

func NewSqlLogger(l ILogger) *SqlLogger {
	return &SqlLogger{logger: l.(*loggerData), showSql: true, level: core.LOG_ERR}
}

func (sl *SqlLogger) Debug(v ...interface{}) {
	if core.LOG_DEBUG >= sl.level {
		sl.logger.writeBuffer(DEBUG_N, 4, v...)
	}
}

func (sl *SqlLogger) Debugf(format string, v ...interface{}) {
	if core.LOG_DEBUG >= sl.level {
		sl.logger.writeBufferf(DEBUG_N, 4, format, v...)
	}
}

func (sl *SqlLogger) Error(v ...interface{}) {
	if core.LOG_ERR >= sl.level {
		sl.logger.writeBuffer(ERROR_N, 4, v...)
	}
}

func (sl *SqlLogger) Errorf(format string, v ...interface{}) {
	if core.LOG_ERR >= sl.level {
		sl.logger.writeBufferf(ERROR_N, 4, format, v...)
	}
}

func (sl *SqlLogger) Info(v ...interface{}) {
	if core.LOG_INFO >= sl.level {
		sl.logger.writeBuffer(INFO_N, 4, v...)
	}
}

func (sl *SqlLogger) Infof(format string, v ...interface{}) {
	if core.LOG_INFO >= sl.level {
		sl.logger.writeBufferf(INFO_N, 4, format, v...)
	}
}

func (sl *SqlLogger) Warn(v ...interface{}) {
	if core.LOG_WARNING >= sl.level {
		sl.logger.writeBuffer(WARN_N, 4, v...)
	}
}

func (sl *SqlLogger) Warnf(format string, v ...interface{}) {
	if core.LOG_WARNING >= sl.level {
		sl.logger.writeBufferf(WARN_N, 4, format, v...)
	}
}

func (sl *SqlLogger) Level() core.LogLevel {
	return sl.level
}

func (sl *SqlLogger) SetLevel(level core.LogLevel) {
	sl.level = level
}

func (sl *SqlLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		sl.showSql = show[0]
	} else {
		sl.showSql = true
	}
}

func (sl *SqlLogger) IsShowSQL() bool {
	return sl.showSql
}
