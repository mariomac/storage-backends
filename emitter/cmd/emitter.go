package main

import (
	"time"

	"github.com/mariomac/storage-backends/emitter/pkg/flow"

	"github.com/mariomac/storage-backends/emitter/pkg/loki"
)

func main() {
	rndGen := flow.NewRndGenerator(4, 20)
	_ = rndGen
	cl := loki.NewHttpJsonClient("http://localhost:3100")
	for {
		err := cl.Push(map[string]string{"source": "fluentd"},
			loki.LogEntry{
				EpochNs: time.Now().UnixNano(),
				Line:    rndGen.Rnd(),
			})
		if err != nil {
			panic(err)
		}
	}
}
