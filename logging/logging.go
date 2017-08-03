package logging

import (
	"errors"
	"log"
	"os"
)

// Logger is struct for logging.
type Logger struct {
	tag string
}

// NewLogger creates new Logger.
func NewLogger(tag string) *Logger {
	return &Logger{tag: tag}
}

var logLv = INFO // default
var lgr = log.New(os.Stderr, "", log.LstdFlags)

// SetLogLevel sets log level to Logger.
func SetLogLevel(leverString string) error {
	var lv level
	switch leverString {
	case "TRACE":
		lv = TRACE
	case "DEBUG":
		lv = DEBUG
	case "INFO":
		lv = INFO
	case "WARNING":
		lv = WARNING
	case "ERROR":
		lv = ERROR
	case "FATAL":
		lv = FATAL
	default:
		return errors.New("logLevel is not valid")
	}
	if logLv != lv {
		logLv = lv
		lgr.SetFlags(log.LstdFlags)
	}

	return nil
}

func (logger *Logger) message(lv level, m string) string {
	return lv.String() + " -- : [" + logger.tag + "] " + m
}

func (logger *Logger) log(lv level, message string, args ...interface{}) {
	if lv >= logLv {
		lgr.Output(3, logger.message(lv, message))
	}
}

// Fatal outputs FATAL log.
func (logger *Logger) Fatal(m string) {
	logger.log(FATAL, m)
}

// Error outputs ERROR log.
func (logger *Logger) Error(m string) {
	logger.log(ERROR, m)
}

// Info outputs INFO log.
func (logger *Logger) Info(m string) {
	logger.log(INFO, m)
}

// Debug outputs DEBUG log.
func (logger *Logger) Debug(m string) {
	logger.log(DEBUG, m)
}

// Trace outputs TRACE log.
func (logger *Logger) Trace(m string) {
	logger.log(TRACE, m)
}
