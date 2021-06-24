package logz

import "log"

var DefaultLogger = NewStdLogger(log.Writer())

type Logger interface {
	Print(level level, vals ...interface{}) error
}

type logger struct {
	logs []Logger
}

func (l *logger) Print(level level, vals ... interface{}) error {
	for _, lg := range l.logs {
		if err := lg.Print(level, vals...); err != nil {
			return err
		}
	}

	return nil
}