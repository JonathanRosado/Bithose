package connectionstore

import "time"

type Message struct {
	LabelPairs []LabelPair `json:"label_pairs"`
	Timestamp  time.Time   `json:"timestamp"`
	Body       string      `json:"body"`
}

type LabelPair struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
