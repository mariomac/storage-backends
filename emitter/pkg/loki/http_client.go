package loki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pathPush    = "/loki/api/v1/push"
	contentPush = "application/json"
)

type HttpJsonClient struct {
	client   *http.Client
	hostPort string
}

type pushPayload struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values []LogEntry        `json:"values"`
}

type LogEntry struct {
	EpochNs int64
	Line    string
}

func (le *LogEntry) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`["%d",%q]`, le.EpochNs, le.Line)), nil
}

func NewHttpJsonClient(hostPort string) HttpJsonClient {
	return HttpJsonClient{
		client:   &http.Client{},
		hostPort: hostPort,
	}
}

func (c *HttpJsonClient) Push(labels map[string]string, entries ...LogEntry) error {
	pl := pushPayload{
		Streams: []stream{{
			Stream: labels,
			Values: entries,
		}},
	}

	payload, _ := json.Marshal(pl)
	_, err := c.client.Post(c.hostPort+pathPush, contentPush, bytes.NewReader(payload))
	return err
}

func (c *HttpJsonClient) Query(labels map[string]string) ([]LogEntry, error) {

}
