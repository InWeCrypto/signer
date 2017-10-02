package logrus

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/goany/slf4go"
)

func TestLogRus(t *testing.T) {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	slf4go.Backend(NewLoggerFactory(log))

	logger := slf4go.Get("test")

	logger.Debug("######")
	logger.Info("????????????")
}
