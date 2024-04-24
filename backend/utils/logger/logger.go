package logger

import (
	"io"
	"os"

	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(debug bool, log *lumberjack.Logger) error {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	var output io.Writer
	if log != nil {
		output = log
	} else {
		output = os.Stdout
	}
	slog.SetDefault(slog.New(newTraceHandler(output, false, level)))
	slog.Info("init log", "output", func() string {
		if log != nil {
			return log.Filename
		} else {
			return "#stdout"
		}
	}())
	return nil
}
