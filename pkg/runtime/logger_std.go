package runtime

import "log"

type StdLogger struct{}

func NewStdLogger() *StdLogger {
	return &StdLogger{}
}

func (l *StdLogger) Infof(format string, args ...any) {
	log.Printf("[INFO] "+format, args...)
}

func (l *StdLogger) Warnf(format string, args ...any) {
	log.Printf("[WARN] "+format, args...)
}

func (l *StdLogger) Errorf(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
}
