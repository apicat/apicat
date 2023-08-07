package log

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
		[]byte(strings.ToUpper(config.SysConfig.Log.Level)),
	); err != nil {
		level = slog.LevelInfo
	}

	if config.SysConfig.Log.Path != "" {
		output = &lumberjack.Logger{
			Filename:   path.Join(config.SysConfig.Log.Path, "apicat.log"),
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
