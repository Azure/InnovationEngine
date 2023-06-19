package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Level string

const (
	Trace Level = "trace"
	Debug Level = "debug"
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
	Fatal Level = "fatal"
)

// / Convert a logging level to a logrus level (uint32).
func (l Level) Integer() logrus.Level {
	switch l {
	case Trace:
		return logrus.TraceLevel
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warn:
		return logrus.WarnLevel
	case Error:
		return logrus.ErrorLevel
	case Fatal:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

// / Convert a string to a logging level.
func LevelFromString(level string) Level {
	switch level {
	case string(Trace):
		return Trace
	case string(Debug):
		return Debug
	case string(Info):
		return Info
	case string(Warn):
		return Warn
	case string(Error):
		return Error
	case string(Fatal):
		return Fatal
	default:
		return Info
	}
}

var GlobalLogger = logrus.New()

func Init(level Level) {
	GlobalLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		DisableQuote:  true,
	})

	GlobalLogger.SetReportCaller(true)
	GlobalLogger.SetLevel(level.Integer())

	file, err := os.OpenFile("ie.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		GlobalLogger.SetOutput(file)
	} else {
		GlobalLogger.SetOutput(os.Stdout)
		GlobalLogger.Info("Failed to log to file, using default stderr")
	}
}
