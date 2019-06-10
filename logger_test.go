package logger


import (
	"testing"
	"context"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {

	config := NewConfig()
	config.Level = "info"
	config.FileName = "logger_test.log"
	SetupLogger(config)
	//
	Log(context.Background()).Info("test")	
	Log(context.Background()).Info("test",
		zap.String("server addr", "xxx"), zap.String("interval", "123"))
}