package main

import (
	"github.com/go-yaml/yaml"
	"github.com/timoguin/goldap/cmd"
	"go.uber.org/zap"
)

var (
	Version   string
	GitCommit string
	GoVersion string

	PlainLogger  *zap.Logger
	Logger       *zap.SugaredLogger
	LoggerConfig *zap.Config
)

func main() {
	// Set built-time vars in the cmd package so we can access them in commands
	cmd.Version = Version
	cmd.GitCommit = GitCommit
	cmd.GoVersion = GoVersion

	// Configure logging
	if err := yaml.Unmarshal([]byte(DefaultLoggerConfig), &LoggerConfig); err != nil {
		panic(err)
	}
	PlainLogger, err := LoggerConfig.Build()
	if err != nil {
		panic(err)
	}
	defer PlainLogger.Sync()

	// Redirect the standard library's global logger
	zap.RedirectStdLogAt(PlainLogger, zap.DebugLevel)

	// Create SugaredLogger
	Logger = PlainLogger.Sugar()

	// Pass loggers to the cmd package
	cmd.PlainLogger = PlainLogger
	cmd.Logger = Logger
	cmd.LoggerConfig = LoggerConfig

	cmd.Execute()
}

const DefaultLoggerConfig = `---
  level: debug
  development: true
  disableCaller: false
  disableStacktrace: false
  sampling: null
  encoderConfig:
    levelKey: level
    messageKey: msg
    levelEncoder: lowercase
  encoding: json
  outputPaths:
    - stdout
    - /tmp/goldap.log
  errorOutputPaths:
    - stderr
    - /tmp/goldap.error.log
  initialFields: {}
}`
