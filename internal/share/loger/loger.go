package loger

import (
	"github.com/sirupsen/logrus"
)

type customLogger struct {
	defaultField string
	formatter    logrus.Formatter
}

func (l customLogger) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["src"] = l.defaultField
	return l.formatter.Format(entry)
}

// New creates a new logger instance with the given name and optional log level.
// The name is used as a source identifier in log entries.
// If no log level is provided, it defaults to TraceLevel.
func New(name string, logLevel ...logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(customLogger{
		defaultField: name,
		formatter:    logrus.StandardLogger().Formatter,
	})
	if len(logLevel) > 0 {
		log.SetLevel(logLevel[0])
	} else {
		log.SetLevel(logrus.TraceLevel)
	}

	return log
}
