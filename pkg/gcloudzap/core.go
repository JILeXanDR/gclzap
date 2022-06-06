package gcloudzap

import (
	"fmt"

	gcl "cloud.google.com/go/logging"
	"github.com/jonstaryuk/gcloudzap"
	"go.uber.org/zap/zapcore"
)

// A Core implements zapcore.Core and writes entries to a Logger from the
// Google Cloud package.
//
// It's safe for concurrent use by multiple goroutines as long as it's not
// mutated after first use.
type Core struct {
	// Logger is a logging.Logger instance from the Google Cloud Platform Go
	// library.
	Logger GoogleCloudLogger

	// Provide your own mapping of zapcore's Levels to Google's Severities, or
	// use DefaultSeverityMapping. All of the Core's children will default to
	// using this map.
	//
	// This must not be mutated after the Core's first use.
	SeverityMapping map[zapcore.Level]gcl.Severity

	// MinLevel is the minimum level for a log entry to be written.
	MinLevel zapcore.Level

	// fields should be built once and never mutated again.
	fields map[string]interface{}
}

func NewCore(client *gcl.Client, gclLogID string, level zapcore.Level) *gcloudzap.Core {
	return &gcloudzap.Core{
		Logger:          client.Logger(gclLogID),
		SeverityMapping: gcloudzap.DefaultSeverityMapping,
		MinLevel:        level,
	}
}

// Enabled implements zapcore.Core.
func (c *Core) Enabled(l zapcore.Level) bool {
	return l >= c.MinLevel
}

// With implements zapcore.Core.
func (c *Core) With(newFields []zapcore.Field) zapcore.Core {
	return &Core{
		Logger:          c.Logger,
		SeverityMapping: c.SeverityMapping,
		MinLevel:        c.MinLevel,
		fields:          clone(c.fields, newFields),
	}
}

// Check implements zapcore.Core.
func (c *Core) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}

// Write implements zapcore.Core. It writes a log entry to Stackdriver.
//
// The "logger", "msg", "caller", and "stack" fields on the payload are
// populated from their respective values on the zapcore.Entry.
func (c *Core) Write(ze zapcore.Entry, newFields []zapcore.Field) error {
	severity, specified := c.SeverityMapping[ze.Level]
	if !specified {
		severity = gcl.Default
	}

	payload := clone(c.fields, newFields)

	payload["logger"] = ze.LoggerName
	payload["msg"] = ze.Message
	payload["caller"] = ze.Caller.String()
	payload["stack"] = ze.Stack

	// TODO: can it log batch?
	c.Logger.Log(gcl.Entry{
		Timestamp: ze.Time,
		Severity:  severity,
		Payload:   payload,
	})

	return nil
}

// Sync implements zapcore.Core. It flushes the Core's Logger instance.
func (c *Core) Sync() error {
	if err := c.Logger.Flush(); err != nil {
		return fmt.Errorf("gcloudzap: flushing Google Cloud logger: %w", err)
	}
	return nil
}
