package main

import (
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging Logging `yaml:"logging"`
	GCL     GCL     `yaml:"gcl"`
}

type Logging struct {
	Level zapcore.Level `yaml:"level"`
}

type GCL struct {
	ServiceAccountPath string        `yaml:"serviceAccountPath"`
	ProjectID          string        `yaml:"projectId"`
	LogID              string        `yaml:"logId"`
	Level              zapcore.Level `yaml:"level"`
}

func loadConfig(path string) (Config, error) {
	defer func() {
		if err := recover(); err != nil {

		}
	}()

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
