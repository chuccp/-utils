package log

type LogFormatter struct {
	TimestampFormat string
}
var defaultFormat = "2006/01/02 15:04:05.999"
var defaultFormatter = &LogFormatter{TimestampFormat:defaultFormat}