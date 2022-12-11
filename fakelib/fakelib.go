// Package fakelib demonstrates how a library should do logging
// and in logger_test.go, this is controlled
// library will be on error level by default, but test/program can control its levels
package fakelib

import (
	"sync"
	"time"

	"github.com/go-msvc/logger"
)

// a library should define its own logger which will have full package name and level error
// so that the user's program will not see DEBUG/INFO logs from the library
// unless requested using e.g. logger.Named("github.com/aname/arepo/alib").SetLevel(logger.LevelDebug)
var log = logger.New()

func Fake() {
	log.Debugf("Fake on debug")
	log.Infof("Fake on debug")
	log.Errorf("Fake on debug")

	//and it can create context loggers too
	wg := sync.WaitGroup{}
	for i := 1; i < 3; i++ {
		myLog := log.With("i", i)
		wg.Add(1)
		go func(l logger.Logger) {
			myLog.Debugf("Fake on debug")
			time.Sleep(time.Millisecond * 100)
			myLog.Infof("Fake on debug")
			time.Sleep(time.Millisecond * 100)
			myLog.Errorf("Fake on debug")
			time.Sleep(time.Millisecond * 100)
			otherFileFunc(l)
			wg.Done()
		}(myLog)
	}
	wg.Wait()
}
