package logging

import (
	gcl "cloud.google.com/go/logging"
	"github.com/jonstaryuk/gcloudzap"
	"go.uber.org/zap/zapcore"
)

type GoogleCloudLoggingOptions struct {
	CredentialsJSON []byte
	ProjectID       string
	LogID           string
	Level           zapcore.Level
}

func NewGoogleCloudLoggingZapCore(client *gcl.Client, gclLogID string, level zapcore.Level) *gcloudzap.Core {
	return &gcloudzap.Core{
		Logger:          client.Logger(gclLogID),
		SeverityMapping: gcloudzap.DefaultSeverityMapping,
		MinLevel:        level,
	}
}
