package utils

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

func LogRuntimeAttributes() {
	logrus.Infof("go version: %s", runtime.Version())
	logrus.Infof("go os/arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("go num cpu: %d", runtime.NumCPU())
	logrus.Infof("go num goroutine: %d", runtime.NumGoroutine())
}
