package utils

import (
	"io"
	"log"
)

// NullLogger stores the default empty logger to be used.
var NullLogger Logger = &nopLogger{}

// LogLevel represents the code for the log severity level.
type LogLevel int

const (
	// INFO log severity code.
	INFO = iota
	// WARN log severity code.
	WARN
	// ERROR log severity code.
	ERROR
)

// Logger defines a simple logging interface
type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// FileLogger represents a high-level logger supporting multiple log levels.
type FileLogger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

// NewFileLogger creates a new FileLogger that writes in the given io.Writer.
func NewFileLogger(w io.Writer, lvl LogLevel) *FileLogger {
	l := &FileLogger{}
	flag := log.Ldate | log.Ltime | log.Lmicroseconds
	if lvl <= INFO {
		l.info = log.New(w, "INFO: ", flag)
	}
	if lvl <= WARN {
		l.warn = log.New(w, "WARN: ", flag)
	}
	if lvl <= ERROR {
		l.error = log.New(w, "ERR: ", flag)
	}
	return l
}

// Infof writes an info event in the log.
func (f *FileLogger) Infof(format string, args ...interface{}) {
	if f.info == nil {
		return
	}
	f.info.Printf(format, args...)
}

// Warningf writes a warning event in the log.
func (f *FileLogger) Warningf(format string, args ...interface{}) {
	if f.warn == nil {
		return
	}
	f.warn.Printf(format, args...)
}

// Errorf writes an error event in the log.
func (f *FileLogger) Errorf(format string, args ...interface{}) {
	if f.error == nil {
		return
	}
	f.error.Printf(format, args...)
}

type nopLogger struct{}

func (*nopLogger) Infof(format string, args ...interface{}) {

}
func (*nopLogger) Warningf(format string, args ...interface{}) {
}

func (*nopLogger) Errorf(format string, args ...interface{}) {
}

func (*nopLogger) Info(string) {

}
func (*nopLogger) Warning(string) {
}

func (*nopLogger) Error(string) {
}
