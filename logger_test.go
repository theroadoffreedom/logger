package logger


import (
	"testing"
	"context"
)

func TestLogger(t *testing.T) {

	config := NewConfig()
	c.Level = "info"
	c.FileName = "logger_test.log"
	c.FileDir = "/Users/ccongdeng/Project"
	SetupLogger(config)
	
	//
	Log(context.Background()).Info("test")	
	Log(context.Background()).Info("test",
		zap.String("server addr", "xxx"), zap.String("interval", "123"))
}