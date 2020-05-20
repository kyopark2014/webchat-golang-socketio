package logger

import (
	"os"

	"github.com/op/go-logging"
)

var loggingEnable = true

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{time:2006-01-02 15:04:05.999} [%{level:.5s}] [%{module}] %{shortfile} %{shortfunc}(): %{message}`,
)

//var format = logging.MustStringFormatter(
//	`%{time:2006-01-02 15:04:05.999} [%{level:.1s}] %{message}`,
//)

// Secure is an an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
type Secure string

func init() {
	SetupLogger(true, "INFO")
}

// Redacted is called whenever anything is logged using Secure
func (p Secure) Redacted() interface{} {
	return logging.Redact(string(p))
}

// SetupLogger is called in initialization part of this service
func SetupLogger(isEnabled bool, level string) {
	loggingEnable = isEnabled

	backend := logging.NewLogBackend(os.Stdout, "", 0)

	// For messages written to backend we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backendFormatter := logging.NewBackendFormatter(backend, format)

	var lvl logging.Level

	switch level {
	case "ERROR":
		lvl = logging.ERROR
	case "WARNING":
		lvl = logging.WARNING
	case "INFO":
		lvl = logging.INFO
	case "DEBUG":
		lvl = logging.DEBUG
	default:
		lvl = logging.INFO
	}

	// Only errors and more severe messages should be sent to backend1
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(lvl, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)
}

// Logger is to express logging information
type Logger struct {
	*logging.Logger
}

// NewLogger returns an instance of logger
func NewLogger(module string) *Logger {
	l := &Logger{logging.MustGetLogger(module)}
	l.ExtraCalldepth = 1
	return l
}

// D writes debug level log
func (l *Logger) D(format string, v ...interface{}) {
	if loggingEnable {
		l.Debugf(format, v...)
	}
}

// W writes warning level log
func (l *Logger) W(format string, v ...interface{}) {
	if loggingEnable {
		l.Warningf(format, v...)
	}
}

// E writes error level log
func (l *Logger) E(format string, v ...interface{}) {
	if loggingEnable {
		l.Errorf(format, v...)
	}
}

// I writes info level log
func (l *Logger) I(format string, v ...interface{}) {
	if loggingEnable {
		l.Infof(format, v...)
	}
}
