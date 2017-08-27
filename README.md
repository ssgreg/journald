# journald
[![GoDoc](https://godoc.org/github.com/ssgreg/journald-send?status.svg)](https://godoc.org/github.com/ssgreg/journald-send)
[![Build Status](https://travis-ci.org/ssgreg/journald-send.svg?branch=master)](https://travis-ci.org/ssgreg/journald-send)
[![Go Report Status](https://goreportcard.com/badge/github.com/ssgreg/journald)](https://goreportcard.com/report/github.com/ssgreg/journald)

Package `journald` offers Go implementation of systemd Journal's native API for logging. Key features are:

* based on a connect-less socket
* work with messages of any size and type
* client can use any number of separate connections


## Example

Let's look at what the `journald` provides as Go APIs for logging:

```go
package main

import (
    "github.com/ssgreg/journald"
)

func main() {
    journald.Print(journald.PriorityInfo, "Hello World!")
}
```

The JSON representation of the journal entry this generates:

```json
{
    "PRIORITY": "6",
    "MESSAGE":  "Hello World!",
    "_PID":     "3965",
    "_COMM":    "simple",
    "...":      "..."
}
```

The primary reason for using the Journal's native logging APIs is not just plain logs: it is to allow passing additional structured log messages from the program into the journal. This additional log data may the be used to search the journal for, is available for consumption for other programs, and might help the administrator to track down issues beyond what is expressed in the human readable message text. Here's an example how to do that with `journals.Send`:

```go
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
```

This will write a log message to the journal much like the earlier examples. However, this times a few additional, structured fields are attached:

```json
{
    "PRIORITY":     "6",
    "MESSAGE":      "Hello World!",
    "HOME":         "/root",
    "TERM":         "xterm",
    "N_GOROUTINE":  "2",
    "N_CPUS":       "4",
    "TRACE"         [103,111,32,49,46,56,32,116,114,97,99,101,0,0,0,0],
    "_PID":         "4037",
    "_COMM":        "send",
    "...":          "..."
}
```

Our structured message includes seven fields. The first two we passed are well-known fields:

1. `MESSAGE=` is the actual human readable message part of the structured message.
1. `PRIORITY=` is the numeric message priority value as known from BSD syslog formatted as an integer string.

Applications may relatively freely define additional fields as they see fit (we defined four pretty arbitrary ones in our example). A complete list of the currently well-known fields is available [here](https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html).

> Thanks to [http://0pointer.de/blog/](http://0pointer.de/blog/) for the inspiration.
