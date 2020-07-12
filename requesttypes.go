package Bithose

import (
	"Bithose/connectionstore"
)

type IncomingSubscribeRequest struct {
	Type     string                                     `json:"type"`
	Criteria []connectionstore.LabelAcceptanceCriterion `json:"criteria"`
}

type IncomingUnsubscribeRequest struct {
	Type string `json:"type"`
	Uuid string `json:"uuid"`
}

type IncomingMessage struct {
	Type    string                  `json:"type"`
	Message connectionstore.Message `json:"message"`
}
