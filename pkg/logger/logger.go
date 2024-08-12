package logger

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
)

type loggerKey int

const MyLoggerKey loggerKey = 3

func InitLogger() (*zap.SugaredLogger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error in logger initialization: %s", err)
		return nil, err
	}
	myLogger := zapLogger.Sugar()
	return myLogger, nil
}

func GetLoggerFromContext(ctx context.Context) (*zap.SugaredLogger, error) {
	myLogger, ok := ctx.Value(MyLoggerKey).(*zap.SugaredLogger)
	if !ok {
		return nil, fmt.Errorf("no logger in context")
	}
	return myLogger, nil
}
