package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mariomac/storage-backends/emitter/pkg/flow"
	"github.com/mariomac/storage-backends/emitter/pkg/loki"
	"golang.org/x/time/rate"
)

const (
	defaultFlowsPerSecond = 2000
	defaultPods           = 20
	defaultNodes          = 4
	defaultLoki           = "http://localhost:3100"
	maxPayloadSize        = 1024 * 1024
)

func main() {
	cfg := parseConfig()
	rndGen := flow.NewRndGenerator(cfg.pods, cfg.nodes)
	accum := flow.NewAccumulator(&rndGen)
	cl := loki.NewHttpJsonClient(cfg.hostAddress)
	totalFlows := 0
	messages := 0
	start := time.Now()
	rateLimiter := rate.NewLimiter(rate.Limit(cfg.flowsPerSecond), 1)
	for {
		time.Sleep(rateLimiter.Reserve().Delay())
		totalFlows++
		if accum.Receive() < maxPayloadSize {
			continue
		}
		pp := accum.Get()
		go func() {
			if err := cl.Push(pp); err != nil {
				log.Print("ERROR sending data:", err)
			}
		}()
		passedSeconds := time.Now().Sub(start).Seconds()
		messages++
		if messages%100 == 0 {
			log.Printf("%4.1f seconds: %d messages %d flows",
				passedSeconds, messages, totalFlows)
		}
	}
}

type config struct {
	hostAddress    string
	pods           int
	nodes          int
	flowsPerSecond int
}

func parseConfig() config {
	cfg := config{}
	var ok bool
	cfg.hostAddress, ok = os.LookupEnv("LOKI_HOST")
	if !ok {
		cfg.hostAddress = defaultLoki
	}
	cfg.pods = defaultPods
	if pstr, ok := os.LookupEnv("PODS"); ok {
		var err error
		cfg.pods, err = strconv.Atoi(pstr)
		if err != nil {
			log.Printf("wrong pods number: %s. Defaulting to %d", err, defaultPods)
		}
	}
	cfg.nodes = defaultNodes
	if pstr, ok := os.LookupEnv("NODES"); ok {
		var err error
		cfg.nodes, err = strconv.Atoi(pstr)
		if err != nil {
			log.Printf("wrong nodes number: %s. Defaulting to %d", err, defaultNodes)
		}
	}
	cfg.flowsPerSecond = defaultFlowsPerSecond
	if fstr, ok := os.LookupEnv("FLOWS_PER_SECOND"); ok {
		var err error
		cfg.flowsPerSecond, err = strconv.Atoi(fstr)
		if err != nil {
			log.Printf("wrong flowsPerSecond: %s. Defaulting to %d", err, defaultFlowsPerSecond)
		}
	}
	return cfg
}
