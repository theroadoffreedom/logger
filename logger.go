//
//
//	config := logger.NewConfig()
//	logger.Init(config)
//
//	usage:
//	logger.log()
//	...
//
package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	zaplog "github.com/pingcap/log"

	cerror "github.com/theroadoffreedom/cerror"
)

type Config struct {
	// Log level.
	Level string `toml:"level" json:"level"`
	// Log format. one of json, text, or console.
	Format string `toml:"format" json:"format"`
	// Disable automatic timestamps in output.
	DisableTimestamp bool `toml:"disable-timestamp" json:"disable-timestamp"`
	// base file name
	FileName string `toml:"FileName" json:"FileName"`
	// file dir
	FileDir string `toml:"FileDir" json:"FileDir"`
	// max file size
	MaxSize uint `toml:"MaxSize" json:"MaxSize"`
	// rotate
	Rotate bool `toml:"Rotate" json:"Rotate"`
}

func NewConfig() *Config {
	c := &Config{}
	c.Level = defaultLogLevelStr
	c.Format = defaultLogFormat
	c.DisableTimestamp = defaultDisableTimestamp
	c.FileName = getCurrentProcessName()
	c.FileDir = getCurrentDir()
	c.MaxSize = defaultLogMaxSize
	c.Rotate = defaultRotate
	return c
}

func SetupLogger(c *Config) {

	// init file config
	fc := newFileLogConfig(c.Rotate,
		c.MaxSize, fmt.Sprintf("%s/%s", c.FileDir, c.FileName))
	lc := newLogConfig(c.Level, c.Format, fc, c.DisableTimestamp)

	// init zap log core
	err := initZapLogger(lc)
	cerror.MustNil(err)

	// init logger
	err = initLogger(lc)
	cerror.MustNil(err)
}

type ctxKeyType int

const ctxLogKey ctxKeyType = iota

// Logger gets a contextual logger from current context.
// contextual logger will output common fields from context.
func Log(ctx context.Context) *zap.Logger {
	if ctxlogger, ok := ctx.Value(ctxLogKey).(*zap.Logger); ok {
		return ctxlogger
	}
	return zaplog.L()
}

// WithConnID attaches connId to context.
func LogWithConnID(ctx context.Context, connID uint32) context.Context {
	var zaplogger *zap.Logger
	if ctxLogger, ok := ctx.Value(ctxLogKey).(*zap.Logger); ok {
		zaplogger = ctxLogger
	} else {
		zaplogger = zaplog.L()
	}
	return context.WithValue(ctx, ctxLogKey, zaplogger.With(zap.Uint32("conn", connID)))
}

// WithKeyValue attaches key/value to context.
func LogWithKeyValue(ctx context.Context, key, value string) context.Context {
	var zaplogger *zap.Logger
	if ctxLogger, ok := ctx.Value(ctxLogKey).(*zap.Logger); ok {
		zaplogger = ctxLogger
	} else {
		zaplogger = zaplog.L()
	}
	return context.WithValue(ctx, ctxLogKey, zaplogger.With(zap.String(key, value)))
}
