package loki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	pathPush      = "/loki/api/v1/push"
	pathQuery     = "/loki/api/v1/query_range"
	contentPush   = "application/json"
	statusSuccess = "success"
)

type HttpJsonClient struct {
	client   *http.Client
	hostPort string
}

type responseBody struct {
	Status string   `json:"status"`
	Data   struct { // ignoring all the other fields
		Result []struct {
			Values []LogEntry `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func (le *LogEntry) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`["%d",%q]`, le.EpochNs, le.Line)), nil
}

func (le *LogEntry) UnmarshalJSON(b []byte) error {
	var entries [2]string
	if err := json.Unmarshal(b, &entries); err != nil {
		return err
	}
	epoch, err := strconv.ParseInt(entries[0], 10, 64)
	if err != nil {
		return err
	}
	le.EpochNs = epoch
	le.Line = entries[1]
	return nil
}

func NewHttpJsonClient(hostPort string) HttpJsonClient {
	return HttpJsonClient{
		client:   &http.Client{},
		hostPort: hostPort,
	}
}

func (c *HttpJsonClient) Push(payload PushPayload) error {
	jp, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(c.hostPort+pathPush, contentPush, bytes.NewReader(jp))
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response %d (%s): %s",
			resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

// just an incomplete test function
// TODO: return all result's entries
// TODO: search by log content
func (c *HttpJsonClient) QueryRange(windowLength time.Duration, labels map[string]string) ([]LogEntry, error) {
	queryUrl := strings.Builder{}
	queryUrl.WriteString(c.hostPort)
	queryUrl.WriteString(pathQuery)
	queryUrl.WriteString("?query={")
	comma := false
	for key, value := range labels {
		if comma {
			queryUrl.WriteByte(',')
		}
		queryUrl.WriteString(key)
		queryUrl.WriteString(`="`)
		queryUrl.WriteString(value)
		queryUrl.WriteString(`"`)
		comma = true
	}
	queryUrl.WriteString("}&start=")
	queryUrl.WriteString(strconv.Itoa(int(time.Now().Add(-windowLength).UnixNano())))

	resp, err := c.client.Get(queryUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected server status: %d. %s", resp.StatusCode, resp.Status)
	}
	var response responseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if response.Status != statusSuccess {
		return nil, fmt.Errorf("response body returned status %q", response.Status)
	}
	// TODO: return complete log entries with their respective labels
	return response.Data.Result[0].Values, nil
}
