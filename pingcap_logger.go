package logger

import (
	"os"
	"path/filepath"
	"strings"

        zaplog "github.com/pingcap/log"
        log "github.com/sirupsen/logrus"
        "gopkg.in/natefinch/lumberjack.v2"

        // self make module
        cerror "github.com/theroadoffreedom/cerror"
)

const (
	defaultLogTimeFormat = "2006/01/02 15:04:05.000"
	defaultLogMaxSize = 300 // MB
	defaultLogFormat = "text"
	defaultLogLevel = "debug"
	defaultDisableTimestamp = false
	defaultRotate = true
)

// LogConfig serializes log related config in toml/json.
type logConfig struct {
        zaplog.Config
}

// NewLogConfig creates a LogConfig.
func newLogConfig(level, format, fileCfg FileLogConfig, disableTimestamp bool) *logConfig {
         return &logConfig{
                 Config: zaplog.Config{
                         Level:            level,
                         Format:           format,
                         DisableTimestamp: disableTimestamp,
                         File:             fileCfg.FileLogConfig,
                 }
         }
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	cerror.MustNil(err)

	return strings.Replace(dir, "\\", "/", -1)
}

func getCurrentProcessName() {
	full := os.Args[0]
	full = strings.Replace(full, "\\", "/", -1)
	splits := strings.Split(full, "/")
	if len(splits) >= 1 {
		name := splits[len(splits)-1]
		name = strings.TrimSuffix(name, ".exe")
		return name
	}
	return ""
}

func stringToLogLevel(level string) log.Level {
         switch strings.ToLower(level) {
         case "fatal":
                 return log.FatalLevel
         case "error":
                 return log.ErrorLevel
         case "warn", "warning":
                 return log.WarnLevel
         case "debug":
                 return log.DebugLevel
         case "info":
                 return log.InfoLevel
         }
         return defaultLogLevel
}

// textFormatter is for compatibility with ngaut/log
type textFormatter struct {
        DisableTimestamp bool
        EnableEntryOrder bool
}

// Format implements logrus.Formatter
func (f *textFormatter) Format(entry *log.Entry) ([]byte, error) {
        var b *bytes.Buffer
        if entry.Buffer != nil {
                b = entry.Buffer
        } else {
                b = &bytes.Buffer{}
        }

        if !f.DisableTimestamp {
                fmt.Fprintf(b, "%s ", entry.Time.Format(defaultLogTimeFormat))
        }
        if file, ok := entry.Data["file"]; ok {
                fmt.Fprintf(b, "%s:%v:", file, entry.Data["line"])
        }
        fmt.Fprintf(b, " [%s] %s", entry.Level.String(), entry.Message)

        if f.EnableEntryOrder {
                keys := make([]string, 0, len(entry.Data))
                for k := range entry.Data {
                        if k != "file" && k != "line" {
                                keys = append(keys, k)
                        }
                }
                sort.Strings(keys)
                for _, k := range keys {
                        fmt.Fprintf(b, " %v=%v", k, entry.Data[k])
                }
        } else {
                for k, v := range entry.Data {
                        if k != "file" && k != "line" {
                                fmt.Fprintf(b, " %v=%v", k, v)
                        }
                }
        }

        b.WriteByte('\n')

        return b.Bytes(), nil
}

func stringToLogFormatter(format string, disableTimestamp bool) log.Formatter {
         switch strings.ToLower(format) {
         case "text":
                 return &textFormatter{
                         DisableTimestamp: disableTimestamp,
                 }
         default:
                 return &textFormatter{}
         }
 }

// InitZapLogger initializes a zap logger with cfg.
func initZapLogger(cfg *LogConfig) error {
        gl, props, err := zaplog.InitLogger(&cfg.Config)
        if err != nil {
                return errors.Trace(err)
        }
        zaplog.ReplaceGlobals(gl, props)
       return nil
}

// InitLogger
func initLogger(cfg *LogConfig) error {
        log.SetLevel(stringToLogLevel(cfg.Level))

        if cfg.Format == "" {
                cfg.Format = defaultLogFormat
        }
        formatter := stringToLogFormatter(cfg.Format, cfg.DisableTimestamp)
        log.SetFormatter(formatter)

        if len(cfg.File.Filename) != 0 {
                if err := initFileLog(&cfg.File, nil); err != nil {
                        return errors.Trace(err)
                }
        }
        return nil
}


// initFileLog initializes file based logging options.
func initFileLog(cfg *zaplog.FileLogConfig, logger *log.Logger) error {
        if st, err := os.Stat(cfg.Filename); err == nil {
                if st.IsDir() {
                        return errors.New("can't use directory as log file name")
                }
        }
        if cfg.MaxSize == 0 {
                cfg.MaxSize = DefaultLogMaxSize
        }

        // use lumberjack to logrotate
        output := &lumberjack.Logger{
                Filename:   cfg.Filename,
                MaxSize:    int(cfg.MaxSize),
                MaxBackups: int(cfg.MaxBackups),
                MaxAge:     int(cfg.MaxDays),
                LocalTime:  true,
        }

        if logger == nil {
                log.SetOutput(output)
        } else {
                logger.Out = output
        }
        return nil
}
