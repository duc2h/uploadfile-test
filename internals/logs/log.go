package logs

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(logLevel string) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"logs.log"}
	config.DisableStacktrace = true
	config.EncoderConfig.TimeKey = "datetime"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapLogger, err := config.Build()
	if err != nil {
		log.Println("cannot build logger, err:", err.Error())
	}
	return zapLogger
}
