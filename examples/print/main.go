package main

import (
	"github.com/ssgreg/journald"
)

func main() {
	journald.Print(journald.PriorityInfo, "Hello World!")
}
