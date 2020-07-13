package client

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

var conn *Connection

func TestMain(m *testing.M) {
	cmd := exec.Command("../bithose")
	err := cmd.Start()
	if err != nil {
		os.Exit(1)
	}

	time.Sleep(time.Second) // wait for server

	c, _ := Connect("localhost:9483")
	conn = c
	defer conn.Close()
	os.Exit(m.Run())
}

func TestSendSubscribe(t *testing.T) {

}

func TestSendMessage(t *testing.T) {
	received := make(chan struct{})

	go func() {
		for {
			message, err := conn.Listen()
			if err != nil {
				t.Error(err.Error())
			}

			if message.Body.(float64) == 6.0 {
				received <- struct{}{}
			}
		}
	}()

	err := conn.Subscribe().
		Criterion("channel", "==", "wowzers").
		Criterion("bugs", "<", 1).
		Send()

	if err != nil {
		t.Error(err.Error())
	}

	err = conn.Message(6).
		Label("channel", "wowzers").
		Label("bugs", 0).
		Send()

	if err != nil {
		t.Error("error on send message: ", err.Error())
	}
	time.Sleep(time.Second * 1) // finish writing

	select {
	case <-received:
	case <-time.After(time.Second * 1):
		t.Error("did not received message")
	}
}
