package main

import (
	"fmt"
	"hash/fnv"
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
	defaultLoki           = "http://localhost:3100"
	defaultPodsBaseAddr   = "172.10.6.0"
	maxPayloadSize        = 1024 * 1024
)

func main() {
	cfg := parseConfig()
	log.Printf("%#v", cfg)
	rndGen := flow.NewRndGenerator(cfg.podsBaseAddr, cfg.pods)
	accum := flow.NewAccumulator(&rndGen)
	cl := loki.NewHttpJsonClient(cfg.hostAddress)
	totalFlows := 0
	messages := 0
	start := time.Now()
	lastReport := start
	lastReportFlows := 0
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
			now := time.Now()
			flowsSecond := float64(totalFlows-lastReportFlows) / now.Sub(lastReport).Seconds()
			lastReport = now
			lastReportFlows = totalFlows
			log.Printf("%.1f seconds: %d messages %d flows (%.0f flows/second)",
				passedSeconds, messages, totalFlows, flowsSecond)
		}
	}
}

type config struct {
	hostAddress    string
	podsBaseAddr   string
	pods           int
	flowsPerSecond int
}

func parseConfig() config {
	cfg := config{}
	var ok bool
	cfg.hostAddress, ok = os.LookupEnv("LOKI_HOST")
	if !ok {
		cfg.hostAddress = defaultLoki
	}

	cfg.podsBaseAddr = defaultPodsBaseAddr
	if hpods, ok := os.LookupEnv("HASH_PODS_BASE"); ok {
		if hashPods, _ := strconv.ParseBool(hpods); hashPods {
			hn, _ := os.Hostname()
			h := fnv.New32a()
			h.Write([]byte(hn))
			hash := h.Sum32()
			cfg.podsBaseAddr = fmt.Sprintf("%d.%d.%d.%d",
				uint8(hash&0xFF),
				uint8((hash>>8)&0xFF),
				uint8((hash>>16)&0xFF),
				uint8((hash>>24)&0xFF))
		}
	}
	cfg.pods = defaultPods
	if pstr, ok := os.LookupEnv("PODS"); ok {
		var err error
		cfg.pods, err = strconv.Atoi(pstr)
		if err != nil {
			log.Printf("wrong pods number: %s. Defaulting to %d", err, defaultPods)
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
