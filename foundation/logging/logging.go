// Package logging is concerned with standardizing the configuration and usage
// of the zap logging package.
package logging

import (
	"github.com/rsb/failure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(service, version string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service":         service,
		"service-version": version,
	}

	log, err := config.Build()
	if err != nil {
		return nil, failure.ToConfig(err, "[zap.NewProductionConfig()] config.Build failed")
	}

	return log.Sugar(), nil
}
