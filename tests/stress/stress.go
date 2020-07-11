package stress

import (
	"Bithose/connectionstore"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var u = url.URL{
	Scheme: "ws",
	Host:   string("localhost:80"),
	Path:   "/",
}

func subscribe(channel string) string {
	var subscribe = `{ 
  "type": "subscribe",
  "criteria": [
	{ "operator": "==", "label_pair": { "name": "channel", "value": "channel_%s" } },
	{ "operator": "<", "label_pair": { "name": "num_of_chars", "value": 5 } }
  ] 
}
`

	return fmt.Sprintf(subscribe, channel)
}

func message(channel string) string {
	var message = `{
  "type": "message",
  "message": {
	"body": "hello",
	"label_pairs": [
	  { "name": "channel", "value": "channel_%s" },
	  { "name": "num_of_chars", "value": 4 }
	]
  }
}
`
	return fmt.Sprintf(message, channel)
}

type TestAgent struct {
	totalMessagesReceived int
	totalDelayInReceived  int64
	totalMessagesSent     int
	totalSubscriptions    int

	done chan struct{}
	mtx  *sync.Mutex
}

func NewTestAgent() *TestAgent {
	return &TestAgent{
		totalMessagesReceived: 0,
		totalDelayInReceived:  0,
		totalMessagesSent:     0,
		totalSubscriptions:    0,
		mtx:                   &sync.Mutex{},
	}
}

func (t *TestAgent) IncrementTotalMessagesReceived() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.totalMessagesReceived++
}

func (t *TestAgent) IncrementTotalDelayInReceived(delay int64) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.totalDelayInReceived += delay
}

func (t *TestAgent) IncrementTotalMessagesSent() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.totalMessagesSent++
}

func (t *TestAgent) GetAverageDelay() float64 {
	return 5.8
}

func (t *TestAgent) RunTest(connections int) {
	for i := 0; i < connections; i++ {
		// add delay to spread out the websocket connections
		time.Sleep(time.Millisecond * 10)
		fmt.Println("connection #", i)
		i := i
		go func() {
			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
			}
			defer c.Close()

			err = c.WriteMessage(websocket.TextMessage, []byte(subscribe(strconv.Itoa(i))))
			if err != nil {
				log.Println("write:", err)
				return
			}

			go func() {
				for {
					_, message, err := c.ReadMessage()
					if err != nil {
						log.Println("read:", err)
						return
					}

					// determine payload type
					var msg connectionstore.Message
					err = json.Unmarshal(message, &msg)
					if err != nil {
						log.Println(err)
						log.Println("type unmarshal unsuccessful")
						continue
					}

					if err != nil {
						log.Println(err)
						log.Println("message unmarshal unsuccessful")
						continue
					}

					t.IncrementTotalMessagesReceived()
					t.IncrementTotalDelayInReceived(time.Now().Unix() - msg.Timestamp.Unix())
				}
			}()

			ticker := time.NewTicker(time.Millisecond * 500)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					err := c.WriteMessage(websocket.TextMessage, []byte(message(strconv.Itoa(i))))
					if err != nil {
						log.Println("write:", err)
						return
					}
					t.IncrementTotalMessagesSent()
				}
			}
		}()
	}
}

func Run() {
	ta := NewTestAgent()
	ta.RunTest(1000)

	time.Sleep(10 * time.Minute)

	fmt.Println("\n\n=============== DONE ===============")

	fmt.Printf("\nTotal messages sent: %v", ta.totalMessagesSent)
	fmt.Printf("\nTotal messages received: %v", ta.totalMessagesReceived)
	fmt.Printf("\nMiss rate: %v%", (ta.totalMessagesReceived/ta.totalMessagesSent)*100)
	fmt.Printf("\nAverage delay in msg receive: %vs", (float64(ta.totalDelayInReceived)/float64(ta.totalMessagesReceived))*100)
}
