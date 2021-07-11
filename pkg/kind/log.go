package kind

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"sigs.k8s.io/kind/pkg/log"
)

// kindLogger implements the Logger interface
type kindLogger struct {
	Log logr.Logger
}

// Warn meets the Logger interface
func (l kindLogger) Warn(message string) {
	l.Warnf("%s", message)
}

// Warnf meets the Logger interface
func (l kindLogger) Warnf(format string, args ...interface{}) {
	l.Log.Info(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// Error meets the Logger interface
func (l kindLogger) Error(message string) {
	l.Errorf("%s", message)
}

// Errorf meets the Logger interface
func (l kindLogger) Errorf(format string, args ...interface{}) {
	l.Log.Info(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

// V meets the Logger interface
func (l kindLogger) V(level log.Level) log.InfoLogger {
	if level == 0 {
		return kindInfoLogger{l.Log}
	}
	return noopInfoLogger{}
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
	l.Infof("%s", message)
}

// Infof meets the InfoLogger interface
func (l kindInfoLogger) Infof(format string, args ...interface{}) {
	l.Log.Info(strings.TrimSpace(fmt.Sprintf(format, args...)))
}

type noopInfoLogger struct{}

func (l noopInfoLogger) Enabled() bool {
	return false
}
func (l noopInfoLogger) Info(message string)                      {}
func (l noopInfoLogger) Infof(format string, args ...interface{}) {}
