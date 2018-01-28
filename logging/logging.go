package logging

import (
	"fmt"
	"os"
	"time"
)

func Info(msg string) {
	fmt.Fprintln(os.Stdout, time.Now().Format("15:04:05")+" INFO "+msg)
}
func Infof(f string, args ...interface{}) {
	fmt.Fprintln(os.Stdout, time.Now().Format("15:04:05")+" INFO "+fmt.Sprintf(f, args...))
}
