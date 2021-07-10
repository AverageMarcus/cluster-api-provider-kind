package kind

import (
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/kind/pkg/log"
)

type Logger struct {
	Log logr.Logger
}

// Warn meets the Logger interface
func (l Logger) Warn(message string) {
	l.Log.Info(message)
}

// Warnf meets the Logger interface
func (l Logger) Warnf(format string, args ...interface{}) {
	l.Log.Info(fmt.Sprintf(format, args...))
}

// Error meets the Logger interface
func (l Logger) Error(message string) {
	l.Log.Info(message)
}

// Errorf meets the Logger interface
func (l Logger) Errorf(format string, args ...interface{}) {
	l.Log.Info(fmt.Sprintf(format, args...))
}

// V meets the Logger interface
func (l Logger) V(level log.Level) log.InfoLogger {
	return InfoLogger{l.Log}
}

// InfoLogger implements the InfoLogger interface and never logs anything
type InfoLogger struct {
	Log logr.Logger
}

// Enabled meets the InfoLogger interface but always returns false
func (l InfoLogger) Enabled() bool {
	return false
}

// Info meets the InfoLogger interface
func (l InfoLogger) Info(message string) {
	l.Log.Info(message)
}

// Infof meets the InfoLogger interface
func (l InfoLogger) Infof(format string, args ...interface{}) {
	l.Log.Info(fmt.Sprintf(format, args...))
}
