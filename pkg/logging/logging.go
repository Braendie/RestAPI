package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// writerHook is a hook from logrus for log data to stdout and file.
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Fire is a function for writing data to a file and stdout.
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

// Levels just returns levels from hook.
func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// it's a standart logger, which returns in GetLogger function.
var e *logrus.Entry

// Logger is a struct for using different logger's types.
type Logger struct {
	*logrus.Entry
}

// GetLogger returns a link of logger.
func GetLogger() *Logger {
	return &Logger{e}
}

// GetLoggerWithField returns special logger with field.
func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

// init creates logger using logrus library, configures this logger,
// creates a storage folder and file: logs/all.logs,
// creates special hook from logrus for log data to stdout and file.
func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
