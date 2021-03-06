package gorm_logger_wrapper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gopherd/log"
	"gorm.io/gorm/logger"
)

const DefaultCalldepth = 2

type Logger struct {
	logger    *log.Logger
	calldepth int
}

func New(logger *log.Logger, calldepth int) *Logger {
	return &Logger{logger: logger, calldepth: calldepth}
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, format string, a ...any) {
	if l.logger.GetLevel() >= log.LevelInfo {
		l.logger.Print(l.calldepth, log.LevelInfo, fmt.Sprintf(format, a...))
	}
}

func (l *Logger) Warn(crx context.Context, format string, a ...any) {
	if l.logger.GetLevel() >= log.LevelWarn {
		l.logger.Print(l.calldepth, log.LevelWarn, fmt.Sprintf(format, a...))
	}
}

func (l *Logger) Error(ctx context.Context, format string, a ...any) {
	if l.logger.GetLevel() >= log.LevelError {
		l.logger.Print(l.calldepth, log.LevelError, fmt.Sprintf(format, a...))
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	var level = log.LevelDebug
	var recordNotFound bool
	if err != nil {
		recordNotFound = errors.Is(err, logger.ErrRecordNotFound)
		if !recordNotFound {
			level = log.LevelInfo
		}
	}
	if l.logger.GetLevel() < level {
		return
	}
	sql, rowsAffected := fc()
	if err != nil {
		if recordNotFound {
			l.logger.Print(l.calldepth, level, fmt.Sprintf("[%s]: record not found", sql))
		} else {
			l.logger.Print(l.calldepth, level, fmt.Sprintf("[%s]: error=%v", sql, err))
		}
	} else {
		l.logger.Print(l.calldepth, level, fmt.Sprintf("[%s]: affected=%d", sql, rowsAffected))
	}
}
