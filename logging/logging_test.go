package logging

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	var logger = NewLogger("test")
	if logger.tag != "test" {
		t.Errorf("tag should be test but %v", logger.tag)
	}
}

func TestSetLogLevel(t *testing.T) {
	SetLogLevel("ERROR")
	if logLv != ERROR {
		t.Errorf("tag should be tag but %v", logLv.String())
	}
}

func TestLog(t *testing.T) {
	SetLogLevel("TRACE")

	var logger = NewLogger("tag")
	logger.Fatal("Fatal log")
	logger.Error("Error log")
	logger.Info("Info log")
	logger.Debug("Debug log")
	logger.Trace("Trace log")
}
