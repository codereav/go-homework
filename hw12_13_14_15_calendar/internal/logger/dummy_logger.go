package logger

import "github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/app"

type DummyLogger struct {
	app.Logger
}

func (l *DummyLogger) Error(_ string)   {}
func (l *DummyLogger) Warning(_ string) {}
func (l *DummyLogger) Info(_ string)    {}
func (l *DummyLogger) Debug(_ string)   {}
