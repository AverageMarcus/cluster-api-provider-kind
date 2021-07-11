package kind

import (
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/kind/pkg/log"
)

// kindLogger implements the Logger interface
type kindLogger struct {
	Log logr.Logger
}

// Warn meets the Logger interface
func (l kindLogger) Warn(message string) {
	l.Log.Info(message)
}

// Warnf meets the Logger interface
func (l kindLogger) Warnf(format string, args ...interface{}) {
	l.Log.Info(fmt.Sprintf(format, args...))
}

// Error meets the Logger interface
func (l kindLogger) Error(message string) {
	l.Log.Info(message)
}

// Errorf meets the Logger interface
func (l kindLogger) Errorf(format string, args ...interface{}) {
	l.Log.Info(fmt.Sprintf(format, args...))
}

// V meets the Logger interface
func (l kindLogger) V(level log.Level) log.InfoLogger {
	return kindInfoLogger{l.Log}
}

// kindInfoLogger implements the InfoLogger interface
type kindInfoLogger struct {
	Log logr.Logger
}

// Enabled meets the InfoLogger interface but always returns false
func (l kindInfoLogger) Enabled() bool {
	return false
}

// Info meets the InfoLogger interface
func (l kindInfoLogger) Info(message string) {
	l.Log.Info(message)
}

// Infof meets the InfoLogger interface
func (l kindInfoLogger) Infof(format string, args ...interface{}) {
	l.Log.Info(fmt.Sprintf(format, args...))
}
