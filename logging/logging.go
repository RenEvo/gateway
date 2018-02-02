package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/renevo/gateway/env"
)

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelError = "ERROR"
)

const (
	envDebug = "GATEWAY_DEBUG"
)

var isDebug = false

func init() {
	isDebug = env.Bool(envDebug)
}

func Debug(msg string) {
	if !isDebug {
		return
	}
	printLevel(levelDebug, msg)
}

func Debugf(f string, args ...interface{}) {
	if !isDebug {
		return
	}
	printLevel(levelDebug, fmt.Sprintf(f, args...))
}

func Info(msg string) {
	printLevel(levelInfo, msg)
}

func Infof(f string, args ...interface{}) {
	printLevel(levelInfo, fmt.Sprintf(f, args...))
}

func Error(msg string) {
	printLevel(levelError, msg)
}

func Errorf(f string, args ...interface{}) {
	printLevel(levelError, fmt.Sprintf(f, args...))
}

func printLevel(level, msg string) {
	fmt.Fprintln(os.Stdout, time.Now().Format("15:04:05")+" "+level+" "+msg)
}
