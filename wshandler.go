package Bithose

import (
	"encoding/json"
	"github.com/JonathanRosado/Bithose/connectionstore"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var Upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Websocket struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewWebsocket(w http.ResponseWriter, r *http.Request) (*Websocket, error) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &Websocket{
		conn: conn,
		mu:   sync.Mutex{},
	}, nil
}

func (w *Websocket) Send(message []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	err := w.conn.WriteMessage(websocket.TextMessage, message)
	return err
}

func WsHandler(writer http.ResponseWriter, request *http.Request) {
	setCors(writer)

	if request.Method == "OPTIONS" {
		return
	}

	ws, err := NewWebsocket(writer, request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	connectionStore := connectionstore.GetMapStore()

	// all connections will use this channel
	ch := make(chan []byte)

	// all uuids for the connections
	uuids := []string{}
	removeConnections := func() {
		for _, uuid := range uuids {
			connectionStore.RemoveConnection(uuid)
		}
	}

	// goroutine listens for sent messages
	go func() {
		for {
			message := <-ch
			err := ws.Send(message)
			if err != nil {
				log.Println("Error while writing")
				log.Println(err)
				removeConnections()
				return
			}
		}
	}()

	for {
		// read incoming payload
		messageType, p, err := ws.conn.ReadMessage()
		// TODO: find a way to close. first value messageType may help
		log.Println("read message")
		if err != nil {
			log.Println("messageType: ", messageType)
			log.Println("p: ", p)
			log.Println(err)
			break
		}

		// determine payload type
		var typeStub = struct {
			Type string `json:"type"`
		}{}
		err = json.Unmarshal(p, &typeStub)
		if err != nil {
			log.Println(err)
			log.Println("type unmarshal unsuccessful")
			continue
		}

		// handle payload based on type
		switch typeStub.Type {
		case "subscribe":
			incomingSubscribe := IncomingSubscribeRequest{}
			err := json.Unmarshal(p, &incomingSubscribe)
			if err != nil {
				log.Println(err)
				continue
			}

			uuid, err := connectionStore.AddConnection(&connectionstore.Connection{
				Ch:                      ch,
				LabelAcceptanceCriteria: incomingSubscribe.Criteria,
			})

			uuids = append(uuids, uuid)

			// send confirmation
			subscribeResponse := SubscribeResponse{
				Uuid: uuid,
			}
			if err != nil {
				subscribeResponse.Error = err.Error()
			}
			jsonResponse, err := json.Marshal(&subscribeResponse)
			if err != nil {
				log.Println(err)
			}
			err = ws.Send(jsonResponse)
			if err != nil {
				log.Println(err)
			}

		case "unsubscribe":
		case "message":
			incomingMessage := IncomingMessage{}
			err := json.Unmarshal(p, &incomingMessage)
			incomingMessage.Message.Timestamp = time.Now()
			if err != nil {
				log.Println(err)
				log.Println("message unmarshal unsuccessful")
				continue
			}

			numOfSent, numOfTimeout, err := connectionStore.SendMessage(incomingMessage.Message)

			// send confirmation
			messageResponse := SendMessageResponse{
				NumberOfSents:    numOfSent,
				NumberOfTimeouts: numOfTimeout,
				Error:            "",
			}
			if err != nil {
				messageResponse.Error = err.Error()
			}
			jsonResponse, err := json.Marshal(&messageResponse)
			if err != nil {
				log.Println(err)
			}
			err = ws.Send(jsonResponse)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("unknown type")
			continue
		}
	}
}

func setCors(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
