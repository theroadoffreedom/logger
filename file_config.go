package logger

import (
	zaplog "github.com/pingcap/log"
)

// FileLogConfig serializes file log related config in toml/json.
type fileLogConfig struct {
        zaplog.FileLogConfig
}

// newFileLogConfig creates a FileLogConfig.
func newFileLogConfig(rotate bool, maxSize uint, fileName string) fileLogConfig {
        return fileLogConfig{fileLogConfig: zaplog.FileLogConfig{
                LogRotate: rotate,
                MaxSize:   int(maxSize),
                Filename : fileName
        }}
}
