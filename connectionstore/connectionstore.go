package connectionstore

import (
	"errors"
	"sync"
)

type ConnectionStore interface {
	AddConnection(connection *Connection) (addedUuid string, err error)
	RemoveConnection(uuid string)
	GetConnection(uuid string) (connection *Connection, exists bool)
	SendMessage(message Message) (numOfSent int, numOfTimeouts int, err error)
	Stats() *Statistics
}

type Connection struct {
	Ch                      chan []byte
	LabelAcceptanceCriteria []LabelAcceptanceCriterion
}

func NewConnection(ch chan []byte, labelAcceptanceCriteria []LabelAcceptanceCriterion) *Connection {
	return &Connection{
		Ch:                      ch,
		LabelAcceptanceCriteria: labelAcceptanceCriteria,
	}
}

// AcceptsLabels returns true if the given Label Pairs matches all the
// criterions specified in the connection. If the label is
func (c *Connection) AcceptsLabels(pairs []LabelPair) (bool, error) {
	pairs = c.filterLabelPairs(pairs)
	if len(pairs) != len(c.LabelAcceptanceCriteria) {
		return false, nil
	}

	accepts := true

	for _, pair := range pairs {
		meetsCriteria := false
		for _, criterion := range c.LabelAcceptanceCriteria {
			accepts, _, err := criterion.acceptsLabel(pair)
			if err != nil {
				return false, err
			}
			if accepts {
				meetsCriteria = true
			}
		}
		if !meetsCriteria {
			accepts = false
		}
	}

	return accepts, nil
}

// filterLabelPairs filters the label pairs such that only the label pairs with label names
// mentioned in the criteria are present. It excludes labels for which there is no
// criteria.
func (c *Connection) filterLabelPairs(pairs []LabelPair) []LabelPair {
	// Load criterion names into a set for quick lookup
	criterionNames := map[string]struct{}{}
	for _, criterion := range c.LabelAcceptanceCriteria {
		criterionNames[criterion.LabelPair.Name] = struct{}{}
	}

	var filteredPairs []LabelPair
	for _, pair := range pairs {
		if _, ok := criterionNames[pair.Name]; ok {
			filteredPairs = append(filteredPairs, pair)
		}
	}

	return filteredPairs
}

type LabelAcceptanceCriterion struct {
	LabelPair LabelPair `json:"label_pair"`
	Operator  string    `json:"operator"`
}

// acceptsLabel takes in a label Name and label value pair from the pushed message and does two things:
//  - checks whether the label Name is included in the LabelAcceptanceCriterion
//  - if it is, checks whether the label value meets the LabelAcceptanceCriterion
// If either of the aforementioned is false, AcceptsLabel returns false. True otherwise.
// NOTE: This implementation assumes that if two label names are equal, the underlying value will be of
// the same type
func (l *LabelAcceptanceCriterion) acceptsLabel(pair LabelPair) (result bool, labelMismatch bool, err error) {
	if l.LabelPair.Name != pair.Name {
		return false, true, nil
	}

	switch pair.Value.(type) {
	case float64: /* since json does not differentiate between int or float, the marshaller will encode numbers as float64 */
		switch l.Operator {
		case "==":
			return pair.Value.(float64) == l.LabelPair.Value.(float64), false, nil
		case ">":
			return pair.Value.(float64) > l.LabelPair.Value.(float64), false, nil
		case "<":
			return pair.Value.(float64) < l.LabelPair.Value.(float64), false, nil
		case ">=":
			return pair.Value.(float64) >= l.LabelPair.Value.(float64), false, nil
		case "<=":
			return pair.Value.(float64) <= l.LabelPair.Value.(float64), false, nil
		default:
			return false, false, OperatorNotFound
		}
	case string:
		switch l.Operator {
		case "==":
			return pair.Value.(string) == l.LabelPair.Value.(string), false, nil
		default:
			return false, false, OperatorNotFound
		}
	case bool:
		switch l.Operator {
		case "==":
			return pair.Value.(bool) == l.LabelPair.Value.(bool), false, nil
		default:
			return false, false, OperatorNotFound
		}
	}

	return false, false, nil
}

type Statistics struct {
	TotalConnections     int `json:"total_connections"`
	TotalMessagesSent    int `json:"total_messages_sent"`
	TotalMessagesTimeout int `json:"total_messages_timeout"`
	mtx                  *sync.RWMutex
}

func NewStatistics() *Statistics {
	return &Statistics{
		TotalConnections:     0,
		TotalMessagesSent:    0,
		TotalMessagesTimeout: 0,
		mtx:                  &sync.RWMutex{},
	}
}

func (s *Statistics) IncrementConnection() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.TotalConnections++
}

func (s *Statistics) DecrementConnection() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.TotalConnections--
}

func (s *Statistics) IncrementMessageSent() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.TotalMessagesSent++
}

func (s *Statistics) IncrementMessageTimeout() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.TotalMessagesTimeout++
}

var (
	OperatorNotFound = errors.New("Operator not found")
)
