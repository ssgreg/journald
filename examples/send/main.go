package main

import (
	"os"
	"runtime"
	"strconv"

	"github.com/ssgreg/journald"
)

func main() {
	journald.Send("Hello World!", journald.PriorityInfo, map[string]string{
		"HOME":        os.Getenv("HOME"),
		"TERM":        os.Getenv("TERM"),
		"N_GOROUTINE": strconv.Itoa(runtime.NumGoroutine()),
		"N_CPUS":      strconv.Itoa(runtime.NumCPU()),
	})
}
