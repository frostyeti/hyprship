package telemetry

import (
	"io"
	"log/slog"
	"os"

	"github.com/frostyeti/hyprship/apps/api/config"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger initializes and returns a new slog.Logger based on config settings
func InitLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler
	logCfg := cfg.Log

	if logCfg == nil {
		logCfg = &config.LogSettings{
			Level:    "info",
			Exporter: "stdout",
			Format:   "text",
		}
	}

	var level slog.Level
	err := level.UnmarshalText([]byte(logCfg.Level))
	if err != nil {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	exporter := logCfg.Exporter
	if exporter == "stdout" && (cfg.Env == "dev" || cfg.Env == "development") {
		exporter = "tint"
	}

	switch exporter {
	case "tint":
		// tint ignores format, uses coloring with isatty for windows support
		w := os.Stdout
		handler = tint.NewHandler(colorable.NewColorable(w), &tint.Options{
			Level:      level,
			NoColor:    !isatty.IsTerminal(w.Fd()),
			TimeFormat: "15:04:05.000",
		})
	case "file":
		file, err := os.OpenFile(logCfg.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("Failed to open log file", "error", err, "filename", logCfg.Filename)
			handler = createHandler(os.Stdout, logCfg.Format, opts)
		} else {
			handler = createHandler(file, logCfg.Format, opts)
		}
	case "rolling-file":
		ljLogger := &lumberjack.Logger{
			Filename:   logCfg.Filename,
			MaxSize:    logCfg.MaxSize,
			MaxBackups: logCfg.MaxBackups,
			MaxAge:     logCfg.MaxAge,
			Compress:   logCfg.Compress,
		}
		handler = createHandler(ljLogger, logCfg.Format, opts)
	case "stderr":
		handler = createHandler(os.Stderr, logCfg.Format, opts)
	default: // stdout or any other string
		handler = createHandler(os.Stdout, logCfg.Format, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func createHandler(w io.Writer, format string, opts *slog.HandlerOptions) slog.Handler {
	if format == "json" {
		return slog.NewJSONHandler(w, opts)
	}
	return slog.NewTextHandler(w, opts)
}
