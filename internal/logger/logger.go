package logger

import (
	"github.com/sirupsen/logrus"
)

type ctxKey string

const keyFieldLogger ctxKey = "fieldLogger"

func FieldLogger(prefix string, value interface{}) logrus.FieldLogger {
	var logger logrus.FieldLogger
	logger = logrus.New().WithField(prefix, value)
	return logger
}
