package log

import (
	"net/url"

	"go.uber.org/zap"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

// Sync implements zap.Sink. The remaining methods are implemented
// by the embedded *lumberjack.Logger.
func (lumberjackSink) Sync() error {
	return nil
}

// NewCustomLogger implements custom zap logger with log roration using an external package lumberjack v2.
func NewCustomLogger() (*zap.Logger, error) {
	zap.RegisterSink("lumberjack", func(u *url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: &lumberjack.Logger{
				Filename:   u.Opaque,
				MaxSize:    100,  // the maximum size in megabytes of the log file before it gets rotated.
				MaxAge:     7,    // the maximum number of days to retain old log files.
				MaxBackups: 10,   // the maximum number of old log files to retain.
				Compress:   true, // determines if the rotated log files should be compressed using gzip.
			},
		}, nil
	})

	// Encoding format of NewDevelopmentConfig is console while NewProductionConfig is json.
	// We don't collect logs to any centralized platform, so use human-readable log format for now.
	config := zap.NewDevelopmentConfig()

	// Create log output file when operating system is linux
	// if runtime.GOOS == "linux" {
	// config.OutputPaths = append(config.OutputPaths, "/home/ubuntu/infra-organizer/logs/chain-exporter.log")
	// }

	return config.Build()
}
