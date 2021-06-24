package logz

import "strings"

type level int8

const (
	LevelDebug level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func(l level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}

	return ""
}

func (l level) ParseLevel(s string) level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	}

	return LevelInfo
}