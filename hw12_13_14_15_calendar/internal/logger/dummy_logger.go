package logger

type DummyLogger struct {
	Logger
}

func (l *DummyLogger) Error(_ string)   {}
func (l *DummyLogger) Warning(_ string) {}
func (l *DummyLogger) Info(_ string)    {}
func (l *DummyLogger) Debug(_ string)   {}
