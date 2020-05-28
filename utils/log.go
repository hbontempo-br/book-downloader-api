package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ErrLogSetup = LogErr("error setting up logger")

type LogErr string

func (e LogErr) Error() string {
	return string(e)
}

// Based on https://stackoverflow.com/questions/57745017/zap-log-framework-go-initialize-log-once-and-reuse-from-other-go-file-solved
func SetupLog(env string) (*zap.Logger, error) {

	var config zap.Config
	switch env {
	case "dev":
		config = zap.NewDevelopmentConfig()
	case "prod":
		config = zap.NewProductionConfig()
	default:
		config = zap.NewDevelopmentConfig()

	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		return nil, ErrLogSetup
	}
	logger.Sugar()
	logger.Debug("Log setup finished successfully")
	return logger, nil

}
