package config

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config log configuration
type Config struct {
	// Level log level
	// default: debug (supports: error、info、trace、warn、panic、fetal、debug)
	Level string

	// Formatter log formatter
	// default: console (supports: json、console)
	// json: json formatter, print log in json format
	// console: console formatter, print log in console
	Formatter string

	// ConsoleSeparator console separator
	// When the formatter is console, the separator between the fields, default is " | "
	ConsoleSeparator string

	// Options zap options
	// default: caller, stacktrace. you can add more options
	Options []zap.Option

	// GlobalFields global fields
	// default: nil
	// global fields, will be added to every log root fields
	GlobalFields map[string]interface{}

	// WriteSyncer write syncer
	// default: os.Stderr
	Syncer zapcore.WriteSyncer

	// ApplyFields apply fields
	// default: nil
	// apply fields, will be added to every log
	ApplyFields func(ctx context.Context) []zap.Field
}
