package appcontext

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/logging"
	"go.uber.org/zap"
)

const (
	defaultLoggerKey = "defaultLoggerKey"
)

var (
	defaultContext AppContext
)

type AppContext struct {
	context.Context
}

func New(ctx context.Context) AppContext {
	appCtx := AppContext{
		Context: ctx,
	}

	defaultContext = appCtx
	return appCtx
}

func GetContext() AppContext {
	return defaultContext
}

func (c *AppContext) getRawContext() context.Context {
	return c.Context
}

func SetLogger(logger *zap.Logger) context.Context {
	defaultContext.Context = context.WithValue(defaultContext.getRawContext(), defaultLoggerKey, logger)
	return defaultContext
}

func GetLogger() *zap.Logger {
	val := defaultContext.Value(defaultLoggerKey)

	v, ok := val.(*zap.Logger)
	if !ok {
		l := logging.NewZapLogger("unknown")
		l.Error("failed to get default logger")
		return l
	}

	return v
}
