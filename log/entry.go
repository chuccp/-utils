package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

type Entry struct {
	Level  Level
	Buffer *bytes.Buffer
	Config *Config
	out    io.Writer

}

func (entry *Entry) WriteTo() {
	entry.Buffer.WriteTo(entry.out)
}
func (entry *Entry) Log(level Level,now   *time.Time, format string, args ...interface{}) {

	tm := now.Format(entry.Config.formatter.TimestampFormat)
	entry.Buffer.WriteString("time=\"" + tm + "\"")
	switch level {
	case TraceLevel:
		entry.Buffer.WriteString(" level=TRACE ")
	case DebugLevel:
		entry.Buffer.WriteString(" level=DEBUG ")
	case InfoLevel:
		entry.Buffer.WriteString(" level=INFO ")
	case WarnLevel:
		entry.Buffer.WriteString(" level=WARN ")
	case ErrorLevel:
		entry.Buffer.WriteString(" level=ERROR ")
	case FatalLevel:
		entry.Buffer.WriteString(" level=FATAL ")
	case PanicLevel:
		entry.Buffer.WriteString(" level=PANIC ")
	}
	entry.Buffer.WriteString("msg=\"")
	end := len(format)
	if end > 0 {
		last := 0
		index := 0
		vLen := len(args)
		for i := 0; i < end; {
			if format[i] == '{' && (i+1 < end) && format[i+1] == '}' {
				entry.Buffer.WriteString(format[last:i])
				if index < vLen {
					fmt.Fprint(entry.Buffer, args[index])
				} else {
					entry.Buffer.WriteString("{}")
				}
				i = i + 2
				index++
				last = i
			}
			i++
		}
		if last < end {
			entry.Buffer.WriteString(format[last:end])
		}

		for index < vLen {
			fmt.Fprint(entry.Buffer, args[index])
			index++
		}
	} else {
		fmt.Fprint(entry.Buffer, args...)
	}
	entry.Buffer.WriteString("\"\n")

}

var poolEntry = &sync.Pool{
	New: func() interface{} {
		return &Entry{Buffer: new(bytes.Buffer)}
	},
}

func newEntry(Config *Config, out io.Writer) *Entry {
	ele := poolEntry.Get().(*Entry)
	ele.Config = Config
	ele.out = out
	ele.Buffer.Reset()
	return ele
}
func freeEntry(p *Entry) {
	poolEntry.Put(p)
}
