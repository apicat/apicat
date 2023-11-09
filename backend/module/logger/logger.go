package logger

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/apicat/apicat/backend/config"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init() {
	var level slog.Level
	var output io.Writer

	if err := level.UnmarshalText(
		[]byte(strings.ToUpper(config.GetSysConfig().Log.Level.Value)),
	); err != nil {
		level = slog.LevelInfo
	}

	if config.GetSysConfig().Log.Path.Value != "" {
		output = &lumberjack.Logger{
			Filename:   path.Join(config.GetSysConfig().Log.Path.Value, "apicat.log"),
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		}
	} else {
		output = os.Stdout
	}
	slog.SetDefault(slog.New(NewTraceHandler(output, false, level)))

}
