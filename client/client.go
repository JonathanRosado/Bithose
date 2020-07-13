package client

import (
	"encoding/json"
	"errors"
	"github.com/JonathanRosado/Bithose"
	"github.com/JonathanRosado/Bithose/connectionstore"
	"github.com/gorilla/websocket"
	"net/url"
)

type Connection struct {
	u    url.URL
	conn *websocket.Conn
}

func Connect(host string) (*Connection, error) {
	var u = url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   "/",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &Connection{
		u:    u,
		conn: conn,
	}, nil
}

func (c *Connection) Close() {
	c.conn.Close()
}

func (c *Connection) Listen() (*connectionstore.Message, error) {
	for {
		var message connectionstore.Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			return nil, err
		}

		// If the receive a non-message, continue
		if message.Body == nil {
			continue
		}

		return &message, nil
	}
}

type Message struct {
	conn    *websocket.Conn
	message *connectionstore.Message
}

func (c *Connection) Message(body interface{}) *Message {
	return &Message{
		conn: c.conn,
		message: &connectionstore.Message{
			LabelPairs: []connectionstore.LabelPair{},
			Body:       body,
		},
	}
}

func (m *Message) Label(name string, value interface{}) *Message {
	m.message.LabelPairs = append(m.message.LabelPairs, connectionstore.LabelPair{
		Name:  name,
		Value: value,
	})
	return m
}

func (m *Message) Send() error {
	message := Bithose.IncomingMessage{
		Type:    "message",
		Message: *m.message,
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = m.conn.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		return err
	}
	return nil
}

type Subscribe struct {
	conn     *websocket.Conn
	criteria []connectionstore.LabelAcceptanceCriterion
}

func (c *Connection) Subscribe() *Subscribe {
	return &Subscribe{
		conn:     c.conn,
		criteria: []connectionstore.LabelAcceptanceCriterion{},
	}
}

func (s *Subscribe) Criterion(name, operator string, value interface{}) *Subscribe {
	s.criteria = append(s.criteria, connectionstore.LabelAcceptanceCriterion{
		LabelPair: connectionstore.LabelPair{
			Name:  name,
			Value: value,
		},
		Operator: operator,
	})
	return s
}

var ErrUnknownOperator = "unknown operator for criterion"

func (s *Subscribe) Send() error {
	// before sending, let's make sure there are no invalid operators
	for _, c := range s.criteria {
		op := c.Operator
		if op != "<" && op != "<=" && op != ">" && op != ">=" && op != "==" {
			return errors.New(ErrUnknownOperator)
		}
	}

	message := Bithose.IncomingSubscribeRequest{
		Type:     "subscribe",
		Criteria: s.criteria,
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = s.conn.WriteMessage(websocket.TextMessage, jsonMessage)
	if err != nil {
		return err
	}
	return nil
}
