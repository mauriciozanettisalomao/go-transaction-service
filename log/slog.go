// Copyright The Linux Foundation and each contributor to LFX.
// SPDX-License-Identifier: MIT

package log

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

const (
	logLevelDefault = slog.LevelDebug

	json = "json"
	text = "text"

	debug = "debug"
	warn  = "warn"
	info  = "info"
)

var (
	logFormatDefault = func(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
		return slog.NewJSONHandler(w, opts)
	}
)

// InitStructureLogConfig sets the structured log behavior
func InitStructureLogConfig() {

	logOptions := &slog.HandlerOptions{}
	var h slog.Handler

	configurations := map[string]func(){
		"options-logLevel": func() {
			logLevel := viper.GetString("logLevel")
			slog.Info("log config",
				"logLevel", logLevel,
			)
			switch logLevel {
			case debug:
				logOptions.Level = slog.LevelDebug
			case warn:
				logOptions.Level = slog.LevelWarn
			case info:
				logOptions.Level = slog.LevelInfo
			default:
				logOptions.Level = logLevelDefault
			}
		},
		"options-addSource": func() {
			addSource := viper.GetBool("logAddSource")
			slog.Info("log config",
				"logAddSource", viper.GetBool("logAddSource"),
			)
			logOptions.AddSource = addSource
		},
		"handler-logFormat": func() {
			logFormat := viper.GetString("logFormat")
			slog.Info("log config",
				"logFormat", logFormat,
			)
			switch logFormat {
			case json:
				h = slog.NewJSONHandler(os.Stdout, logOptions)
			case text:
				h = slog.NewTextHandler(os.Stdout, logOptions)
			default:
				h = logFormatDefault(os.Stdout, logOptions)
			}
		},
	}

	for name, f := range configurations {
		slog.Info("setting logging configuration",
			"name", name,
		)
		f()
	}
	log.SetFlags(log.Llongfile)
	slog.SetDefault(slog.New(h))

}
