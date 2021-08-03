package main

import (
	"time"

	"github.com/mariomac/storage-backends/emitter/pkg/loki"
)

func main() {
	cl := loki.NewHttpJsonClient("http://localhost:3100")
	err := cl.Push(map[string]string{
		"foo": "bar",
		"baz": "bae",
	},
		loki.LogEntry{
			EpochNs: time.Now().UnixNano(),
			Line:    " blabla blabla",
		},
		loki.LogEntry{
			EpochNs: time.Now().UnixNano(),
			Line:    "12\\\"341325",
		},
		loki.LogEntry{
			EpochNs: time.Now().UnixNano(),
			Line:    " 65565656565",
		},
		loki.LogEntry{
			EpochNs: time.Now().UnixNano(),
			Line:    "l 9090909090",
		})
	if err != nil {
		panic(err)
	}

}
