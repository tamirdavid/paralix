// logger/logger.go

package logger

import (
	"github.com/sirupsen/logrus"
)

var Log = func() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return log
}()
