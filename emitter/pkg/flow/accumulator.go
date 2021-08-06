package flow

import (
	"github.com/mariomac/storage-backends/emitter/pkg/loki"
	"time"
)

type index struct {
	srcPod string
	dstPod string
}

type generator interface {
	Generate() (payload, srcPod, dstPod string)
}

type Accumulator struct {
	generator generator
	payloads map[index][]loki.LogEntry
	clock func() time.Time
	accumulatedSize int
}

func NewAccumulator(generator generator) Accumulator {
	return Accumulator{
		generator: generator,
		payloads: map[index][]loki.LogEntry{},
		clock: time.Now,
	}
}

// Receive generates and accumulates one flow
func (a *Accumulator) Receive() int {
	var idx index
	var payload string
	payload, idx.srcPod, idx.dstPod = a.generator.Generate()
	a.payloads[idx] = append(a.payloads[idx], loki.LogEntry{
		EpochNs: a.clock().UnixNano(),
		Line: payload,
	})
	a.accumulatedSize += len(payload)
	return a.accumulatedSize
}

func (a *Accumulator) Get() loki.PushPayload {
	pp := loki.PushPayload{}
	for keys, lines := range a.getAndResetPayloads() {
		pp.Streams = append(pp.Streams, loki.Stream{
			Stream: map[string]string{
				"srcPod": keys.srcPod,
				"dstPod": keys.dstPod,
			},
			Values: lines,
		})
	}
	return pp
}

func (a *Accumulator) getAndResetPayloads() map[index][]loki.LogEntry {
	pl := a.payloads
	a.payloads = map[index][]loki.LogEntry{}
	a.accumulatedSize = 0
	return pl
}
