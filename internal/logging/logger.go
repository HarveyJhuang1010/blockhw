package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger create zap logger by env
func NewZapLogger(labels ...string) (logger *zap.Logger) {
	var (
		err error
	)

	logger, err = zap.NewDevelopment(zap.AddStacktrace(zapcore.DPanicLevel))
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)

	if len(labels) > 0 {
		logger = logger.Named(labels[0])
	}
	return
}
