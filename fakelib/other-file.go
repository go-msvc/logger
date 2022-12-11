package fakelib

import (
	"time"

	"github.com/go-msvc/logger"
)

func otherFileFunc(myLog logger.Logger) {
	myLog.Debugf("Fake on debug")
	time.Sleep(time.Millisecond * 100)
	myLog = myLog.With("valid", true)
	myLog.Infof("Fake on debug")
	time.Sleep(time.Millisecond * 100)
	myLog.Error("Fake on debug")
	time.Sleep(time.Millisecond * 100)
}
