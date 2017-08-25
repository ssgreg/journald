package main

import (
	"github.com/ssgreg/journald-send"
)

func main() {
	journald.Print(journald.PriorityInfo, "Hello World!")
}
