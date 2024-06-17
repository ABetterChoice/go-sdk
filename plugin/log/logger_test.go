// Package log logger
package log

import (
	"fmt"
	"testing"
)

func TestRegister(t *testing.T) {

	// RegisterLogger(nil)
	RegisterLogger(&InnerLogger{})
	SetLoggerLevel(DebugLevel)
	fmt.Println("test debug loggerLevel")
	Debug("hello debug1")
	Debugf("%v: hello debugf1", "tab")

	Info("hello info1")
	Infof("%v: hello infof1", "tab")

	Warn("hello warn1")
	Warnf("%v: hello warnf1", "tab")

	Error("hello error1")
	Errorf("%v: hello errorf1", "tab")

	SetLoggerLevel(InfoLevel)
	fmt.Println("test info loggerLevel")
	Debug("hello debug2")
	Debugf("%v: hello debugf2", "tab")

	Info("hello info2")
	Infof("%v: hello infof2", "tab")

	Warn("hello warn2")
	Warnf("%v: hello warnf2", "tab")

	Error("hello error2")
	Errorf("%v: hello errorf2", "tab")

	SetLoggerLevel(WarnLevel)
	fmt.Println("test warn loggerLevel")
	Debug("hello debug3")
	Debugf("%v: hello debugf3", "tab")

	Info("hello info3")
	Infof("%v: hello infof3", "tab")

	Warn("hello warn3")
	Warnf("%v: hello warnf3", "tab")

	Error("hello error4")
	Errorf("%v: hello errorf4:%v", "tab", "test")

	SetLoggerLevel(ErrorLevel)
	fmt.Println("test error loggerLevel4")
	Debug("hello debug4")
	Debugf("%v: hello debugf4", "tab")

	Info("hello info4")
	Infof("%v: hello infof4", "tab")

	Warn("hello warn4")
	Warnf("%v: hello warnf4", "tab")

	Error("hello error4")
	Errorf("%v: hello errorf4:%v", "tab", "test")

	SetLoggerLevel(NotLogLevel)
	fmt.Println("test notLog loggerLevel5")
	Debug("hello debug5")
	Debugf("%v: hello debugf5", "tab")

	Info("hello info5")
	Infof("%v: hello infof5", "tab")

	Warn("hello warn5")
	Warnf("%v: hello warnf5", "tab")

	Error("hello error5")
	Errorf("%v: hello errorf5:%v", "tab", "test")
}
