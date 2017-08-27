package main

import (
	"os"
	"runtime"

	"github.com/ssgreg/journald"
)

func main() {
	journald.Send("Hello World!", journald.PriorityInfo, map[string]interface{}{
		"HOME":        os.Getenv("HOME"),
		"TERM":        os.Getenv("TERM"),
		"N_GOROUTINE": runtime.NumGoroutine(),
		"N_CPUS":      runtime.NumCPU(),
		"TRACE":       runtime.ReadTrace(),
	})
}
