package logging

import "fmt"

type level uint8

// define log levels.
const (
	_       level = iota
	TRACE         // TRACE log level
	DEBUG         // DEBUG log level
	INFO          // INFO log level
	WARNING       // WARNING log level
	ERROR         // ERROR log level
	FATAL         // FATAL log level
)

func (l level) String() string {
	switch l {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return fmt.Sprintf("level(%d)", l)
}
