package logger

import "testing"

//fun test| debug test
func TestFileLogger(t *testing.T)  {
	logger := NewFileLogger(LogLevelDebug, "/Users/linger/logs/","test")
	logger.Debug("user id[%d] is come from china", 123)
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Close()

}

func TestConsoleLogger(t *testing.T)  {
	logger := NewConsoleLogger(LogLevelDebug)
	logger.Debug("user id[%d] is come from china", 88123)
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Close()

}