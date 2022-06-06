package logging

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	gcl "cloud.google.com/go/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"
	"logging/pkg/gcloudzap"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Options struct {
	Level zapcore.Level
	GCL   GoogleCloudLoggingOptions
}

type GoogleCloudLoggingOptions struct {
	CredentialsJSON []byte
	ProjectID       string
	LogID           string
	Level           zapcore.Level
}

func New(options Options) (*zap.Logger, error) {
	if options.GCL.Level < options.Level {
		return nil, fmt.Errorf(`cloud logging level must be equal or higher than stdout level "%s", but it's set to "%s"`, options.Level.String(), options.GCL.Level.String())
	}

	gclClient, err := gcl.NewClient(context.Background(), options.GCL.ProjectID, option.WithCredentialsJSON(options.GCL.CredentialsJSON))
	if err != nil {
		return nil, err
	}

	if err := gclClient.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("google cloud logging is not available: %w", err)
	}

	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	baseCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, options.Level)
	gclCore := gcloudzap.NewCore(gclClient, options.GCL.LogID, options.GCL.Level)
	core := zapcore.NewTee(baseCore, gclCore)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)) // zap.AddCallerSkip(1)

	return logger, nil
}
