package logz

import "fmt"

type Log struct {
	logger Logger
}

func NewLog(logger Logger) *Log {
	return &Log{
		logger: logger,
	}
}

func (l *Log) Print(level level, vals ...interface{}) {
	l.logger.Print(level, vals...)
}

func (l *Log) Debug(v ...interface{}) {
	l.logger.Print(LevelDebug, fmt.Sprint(v...))
}

func (l *Log) Debugf(format string, v ...interface{}) {
	l.logger.Print(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Log) Debugw(vals ...interface{}) {
	l.logger.Print(LevelDebug, vals...)
}

func (l *Log) Info(v ...interface{}) {
	l.logger.Print(LevelInfo, fmt.Sprint(v...))
}

func (l *Log) Infof(format string, v ...interface{}) {
	l.logger.Print(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Log) Infow(vals ...interface{}) {
	l.logger.Print(LevelInfo, vals...)
}

func (l *Log) Warn(v ...interface{}) {
	l.logger.Print(LevelWarn, fmt.Sprint(v...))
}

func (l *Log) Warnf(format string, v ...interface{}) {
	l.logger.Print(LevelWarn, fmt.Sprintf(format, v...))
}

func (l *Log) Warnw(vals ...interface{}) {
	l.logger.Print(LevelWarn, vals...)
}

func (l *Log) Error(v ...interface{}) {
	l.logger.Print(LevelError, fmt.Sprint(v...))
}

func (l *Log) Errorf(format string, v ...interface{}) {
	l.logger.Print(LevelError, fmt.Sprintf(format, v...))
}

func (l *Log) Errorw(vals ...interface{}) {
	l.logger.Print(LevelError, vals...)
}

