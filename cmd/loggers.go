package cmd

import (
	"go.uber.org/zap"
)

var (
	PlainLogger  *zap.Logger
	Logger       *zap.SugaredLogger
	LoggerConfig *zap.Config
)
