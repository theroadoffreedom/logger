package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {

	config := NewConfig()
	config.Level = "debug"
	config.FileName = "logger_test.log"
	config.FileDir = "/tmp"
	SetupLogger(config)
	//
	Log(context.Background()).Info("test")
	Log(context.Background()).Debug("test",
		zap.String("server addr", "xxx"), zap.String("interval", "123"))
}
