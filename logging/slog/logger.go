package slog

import (
	"fmt"
	stdslog "log/slog"

	"github.com/matiasmartin-labs/common-fwk/logging"
)

type loggerAdapter struct {
	base *stdslog.Logger
}

func newLoggerAdapter(base *stdslog.Logger) logging.Logger {
	return &loggerAdapter{base: base}
}

func (l *loggerAdapter) Debugf(format string, args ...any) {
	l.base.Debug(fmt.Sprintf(format, args...))
}

func (l *loggerAdapter) Infof(format string, args ...any) {
	l.base.Info(fmt.Sprintf(format, args...))
}

func (l *loggerAdapter) Warnf(format string, args ...any) {
	l.base.Warn(fmt.Sprintf(format, args...))
}

func (l *loggerAdapter) Errorf(format string, args ...any) {
	l.base.Error(fmt.Sprintf(format, args...))
}
