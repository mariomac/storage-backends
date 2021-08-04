package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mariomac/storage-backends/emitter/pkg/flow"

	"github.com/mariomac/storage-backends/emitter/pkg/loki"
)

const (
	defaultPods  = 20
	defaultNodes = 4
	defaultLoki  = "http://localhost:3100"
)

func main() {
	hostAddress, pods, nodes := parseConfig()
	rndGen := flow.NewRndGenerator(pods, nodes)
	cl := loki.NewHttpJsonClient(hostAddress)
	start := time.Now()
	messages := 0
	for {
		err := cl.Push(map[string]string{"source": "fluentd"},
			loki.LogEntry{
				EpochNs: time.Now().UnixNano(),
				Line:    rndGen.Rnd(),
			})
		if err != nil {
			panic(err)
		}
		messages++
		if messages%10_000 == 0 {
			log.Printf("%4.1f seconds: %d messages",
				time.Now().Sub(start).Seconds(), messages)
		}
	}
}

func parseConfig() (string, int, int) {
	hostAddress, ok := os.LookupEnv("LOKI_HOST")
	if !ok {
		hostAddress = defaultLoki
	}
	pods := defaultPods
	if pstr, ok := os.LookupEnv("PODS"); ok {
		var err error
		pods, err = strconv.Atoi(pstr)
		if err != nil {
			log.Printf("wrong pods number: %s. Defaulting to %d", err, defaultPods)
		}
	}
	nodes := defaultNodes
	if pstr, ok := os.LookupEnv("NODES"); ok {
		var err error
		nodes, err = strconv.Atoi(pstr)
		if err != nil {
			log.Printf("wrong nodes number: %s. Defaulting to %d", err, defaultNodes)
		}
	}
	return hostAddress, pods, nodes
}
