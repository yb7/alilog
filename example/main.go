package main

import (
	"time"

	"github.com/yb7/alilog"
)

func main() {
	aliLog := alilog.New("cls-log-test", "test-log")
	i := 0
	for i < 1000 {
		i++
		aliLog.Debugf("log %d", i)
	}

	time.Sleep(5 * time.Second)
}
