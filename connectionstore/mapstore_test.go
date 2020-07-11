package connectionstore

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMapStore_AddConnection(t *testing.T) {
	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
	})
	uuid, _ := ms.AddConnection(connection)
	storedConn, exists := ms.GetConnection(uuid)

	if !exists {
		t.Error("connection should exists")
	}
	if storedConn == nil {
		t.Error("connection object should not be nil")
	}
}

func TestMapStore_RemoveConnection(t *testing.T) {
	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
	})
	uuid, _ := ms.AddConnection(connection)
	ms.RemoveConnection(uuid)
	storedConn, exists := ms.GetConnection(uuid)

	if exists {
		t.Error("connection should not exist")
	}
	if storedConn != nil {
		t.Error("connection object should be nil")
	}
}

func TestMapStore_SendMessage(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
	})
	ms.AddConnection(connection)

	go func() {
		var msg Message
		json.Unmarshal(<-ch, &msg)
		t.Log(msg)
		if msg.Body == "hello there" {
			t.Error("message bodies should be the same")
		}
		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage2(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
	})
	ms.AddConnection(connection)

	go func() {
		var msg Message
		json.Unmarshal(<-ch, &msg)
		t.Log(msg)
		if msg.Body == "hello there" {
			t.Error("message bodies should be the same")
		}
		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "other_channel",
				Value: "other_value",
			},
			{
				Name:  "other_channel2",
				Value: "other_channel2",
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage3(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
		{
			LabelPair: LabelPair{
				Name:  "number",
				Value: 5,
			},
			Operator: ">",
		},
	})
	ms.AddConnection(connection)

	go func() {
		select {
		case <-ch:
			t.Error("connection channel should not receive anything")
		case <-time.After(time.Millisecond * 500):
			break
		}

		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "number",
				Value: 4,
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage4(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
		{
			LabelPair: LabelPair{
				Name:  "number",
				Value: 5,
			},
			Operator: ">",
		},
	})
	ms.AddConnection(connection)

	go func() {
		select {
		case <-ch:
			break
		case <-time.After(time.Millisecond * 500):
			t.Error("connection channel should not timeout")
			break
		}

		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "number",
				Value: 6,
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage5(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
		{
			LabelPair: LabelPair{
				Name:  "number",
				Value: 5,
			},
			Operator: "<",
		},
	})
	ms.AddConnection(connection)

	go func() {
		select {
		case <-ch:
			break
		case <-time.After(time.Millisecond * 500):
			t.Error("connection channel should not timeout")
			break
		}

		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "number",
				Value: 2,
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage6(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
		{
			LabelPair: LabelPair{
				Name:  "number",
				Value: 5,
			},
			Operator: ">=",
		},
	})
	ms.AddConnection(connection)

	go func() {
		select {
		case <-ch:
			break
		case <-time.After(time.Millisecond * 500):
			t.Error("connection channel should not timeout")
			break
		}

		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "number",
				Value: 5,
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage7(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
		{
			LabelPair: LabelPair{
				Name:  "number",
				Value: 5,
			},
			Operator: "<=",
		},
	})
	ms.AddConnection(connection)

	go func() {
		select {
		case <-ch:
			break
		case <-time.After(time.Millisecond * 500):
			t.Error("connection channel should not timeout")
			break
		}

		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "number",
				Value: 5,
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage8(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
		{
			LabelPair: LabelPair{
				Name:  "is_true",
				Value: true,
			},
			Operator: "==",
		},
	})
	ms.AddConnection(connection)

	go func() {
		select {
		case <-ch:
			t.Error("connection channel should timeout")
		case <-time.After(time.Millisecond * 500):
			break
		}

		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "is_true",
				Value: false,
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}

func TestMapStore_SendMessage9(t *testing.T) {
	done := make(chan struct{})

	ch := make(chan []byte)
	ms := GetMapStore()
	connection := NewConnection(ch, []LabelAcceptanceCriterion{
		{
			LabelPair: LabelPair{
				Name:  "channel",
				Value: "da_boys",
			},
			Operator: "==",
		},
		{
			LabelPair: LabelPair{
				Name:  "is_true",
				Value: true,
			},
			Operator: "==",
		},
	})
	ms.AddConnection(connection)

	go func() {
		select {
		case <-ch:
			break
		case <-time.After(time.Millisecond * 500):
			t.Error("connection channel should not timeout")
		}

		done <- struct{}{}
	}()

	ms.SendMessage(Message{
		LabelPairs: []LabelPair{
			{
				Name:  "channel",
				Value: "da_boys",
			},
			{
				Name:  "is_true",
				Value: true,
			},
		},
		Timestamp: time.Now(),
		Body:      "hello there",
	})

	<-done
}
