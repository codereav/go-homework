package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

type Logger struct {
	logger *logrus.Logger
}

func New(level string, path string) *Logger {
	if path == "" {
		path = "./"
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic(errors.Wrap(err, "fail on log level parsing"))
	}

	logger := logrus.New()
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(errors.Wrap(err, "fail on log dir creating"))
	}

	var logFile io.Writer
	logfilePath := fmt.Sprintf("%slogfile-%s.log", path, time.DateOnly)
	logFile, err = os.OpenFile(logfilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Info(errors.Wrap(err, "fail on log file opening, will use standard output"))
	} else {
		logger.SetOutput(logFile)
		if logLevel == logrus.DebugLevel { // Also it will write to terminal if level is debug
			hook := &writer.Hook{
				Writer:    os.Stderr,
				LogLevels: logrus.AllLevels,
			}
			logger.Hooks.Add(hook)
		}
	}

	logger.SetLevel(logLevel)

	return &Logger{
		logger: logger,
	}
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Warning(msg string) {
	l.logger.Warning(msg)
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
