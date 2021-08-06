package loki

type PushPayload struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Stream map[string]string `json:"stream"`
	Values []LogEntry        `json:"values"`
}

type LogEntry struct {
	EpochNs int64
	Line    string
}
