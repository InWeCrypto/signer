package slf4go

import "testing"

func TestLog(t *testing.T) {
	logger := Get("test")

	logger.Debug("test")
	logger.DebugF("test %s", "hello")
	logger.Error("error")
	logger.ErrorF("error %p", t)
}
