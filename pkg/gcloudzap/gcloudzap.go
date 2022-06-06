/*
Package gcloudzap provides a zap logger that forwards entries to the Google
Stackdriver Logging service as structured payloads.

All zap.Logger instances created with this package are safe for concurrent
use.

Network calls (which are delegated to the Google Cloud Platform package) are
asynchronous and payloads are buffered. These benchmarks, on a MacBook Pro 2.4
GHz Core i5, are a loose approximation of latencies on the critical path for
the zapcore.Core implementation provided by this package.

	$ go test -bench . github.com/jonstaryuk/gcloudzap
	goos: darwin
	goarch: amd64
	pkg: github.com/jonstaryuk/gcloudzap
	BenchmarkCoreClone-4   	 2000000	       607 ns/op
	BenchmarkCoreWrite-4   	 1000000	      2811 ns/op


Zap docs: https://godoc.org/go.uber.org/zap

Stackdriver Logging docs: https://cloud.google.com/logging/docs/

*/
package gcloudzap

import (
	"fmt"
	"time"

	gcl "cloud.google.com/go/logging"
	"go.uber.org/zap/zapcore"
)

// DefaultSeverityMapping is the default mapping of zap's Levels to Google's
// Severities.
var DefaultSeverityMapping = map[zapcore.Level]gcl.Severity{
	zapcore.DebugLevel:  gcl.Debug,
	zapcore.InfoLevel:   gcl.Info,
	zapcore.WarnLevel:   gcl.Warning,
	zapcore.ErrorLevel:  gcl.Error,
	zapcore.DPanicLevel: gcl.Critical,
	zapcore.PanicLevel:  gcl.Critical,
	zapcore.FatalLevel:  gcl.Critical,
}

// clone creates a new field map without mutating the original.
func clone(orig map[string]interface{}, newFields []zapcore.Field) map[string]interface{} {
	clone := make(map[string]interface{})

	for k, v := range orig {
		clone[k] = v
	}

	for _, f := range newFields {
		switch f.Type {
		// case zapcore.UnknownType:
		case zapcore.ArrayMarshalerType:
			clone[f.Key] = f.Interface
		case zapcore.ObjectMarshalerType:
			clone[f.Key] = f.Interface
		case zapcore.BinaryType:
			clone[f.Key] = f.Interface
		case zapcore.BoolType:
			clone[f.Key] = (f.Integer == 1)
		case zapcore.ByteStringType:
			clone[f.Key] = f.String
		case zapcore.Complex128Type:
			clone[f.Key] = fmt.Sprint(f.Interface)
		case zapcore.Complex64Type:
			clone[f.Key] = fmt.Sprint(f.Interface)
		case zapcore.DurationType:
			clone[f.Key] = time.Duration(f.Integer).String()
		case zapcore.Float64Type:
			clone[f.Key] = float64(f.Integer)
		case zapcore.Float32Type:
			clone[f.Key] = float32(f.Integer)
		case zapcore.Int64Type:
			clone[f.Key] = int64(f.Integer)
		case zapcore.Int32Type:
			clone[f.Key] = int32(f.Integer)
		case zapcore.Int16Type:
			clone[f.Key] = int16(f.Integer)
		case zapcore.Int8Type:
			clone[f.Key] = int8(f.Integer)
		case zapcore.StringType:
			clone[f.Key] = f.String
		case zapcore.TimeType:
			clone[f.Key] = f.Interface.(time.Time)
		case zapcore.Uint64Type:
			clone[f.Key] = uint64(f.Integer)
		case zapcore.Uint32Type:
			clone[f.Key] = uint32(f.Integer)
		case zapcore.Uint16Type:
			clone[f.Key] = uint16(f.Integer)
		case zapcore.Uint8Type:
			clone[f.Key] = uint8(f.Integer)
		case zapcore.UintptrType:
			clone[f.Key] = uintptr(f.Integer)
		case zapcore.ReflectType:
			clone[f.Key] = f.Interface
		// case zapcore.NamespaceType:
		case zapcore.StringerType:
			clone[f.Key] = f.Interface.(fmt.Stringer).String()
		case zapcore.ErrorType:
			clone[f.Key] = f.Interface.(error).Error()
		case zapcore.SkipType:
			continue
		default:
			clone[f.Key] = f.Interface
		}
	}

	return clone
}
