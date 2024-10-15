package logging

import (
	"fmt"
	"log/slog"
)

type RestyLogger struct {
	logger *slog.Logger
}

func NewRestyLogger(logger *slog.Logger) *RestyLogger {
	return &RestyLogger{
		logger: logger,
	}
}

func (r *RestyLogger) Errorf(format string, v ...any) {
	r.logger.Error(fmt.Sprintf(format, v...))
}

func (r *RestyLogger) Warnf(format string, v ...any) {
	r.logger.Warn(fmt.Sprintf(format, v...))
}

func (r *RestyLogger) Infof(format string, v ...any) {
	r.logger.Info(fmt.Sprintf(format, v...))
}

func (r *RestyLogger) Debugf(format string, v ...any) {
	r.logger.Debug(fmt.Sprintf(format, v...))
}
