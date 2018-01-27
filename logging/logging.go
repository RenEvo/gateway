package logging

import (
	"fmt"
	"time"
)

func Info(msg string) {
	fmt.Println(time.Now().Format("15:04:05") + " INFO " + msg)
}
func Infof(f string, args ...interface{}) {
	fmt.Println(time.Now().Format("15:04:05") + " INFO " + fmt.Sprintf(f, args...))
}
