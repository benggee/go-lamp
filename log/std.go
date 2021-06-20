package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log *log.Logger
	pool *sync.Pool
}


func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (l *stdLogger) Print(level level, vals ...interface{}) error {
	buf := l.pool.Get().(*bytes.Buffer)
	buf.WriteString(fmt.Sprintf("%s [%s] ", time.Now().Format("2006-01-02 15:04:05"), level.String()))

	for i := 0; i < len(vals); i ++ {
		fmt.Fprint(buf, vals...)
	}

	l.log.Output(4, buf.String())
	buf.Reset()
	l.pool.Put(buf)
	return nil
}