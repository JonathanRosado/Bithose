package Bithose

import (
	"Bithose/connectionstore"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var Upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(writer http.ResponseWriter, request *http.Request) {
	setCors(writer)

	if request.Method == "OPTIONS" {
		return
	}

	conn, err := Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
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
			err := conn.WriteMessage(websocket.TextMessage, message)
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
		_, p, err := conn.ReadMessage()
		log.Println("read message")
		if err != nil {
			log.Println(err)
			log.Println("closed on read")
			removeConnections()
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

			uuid, _ := connectionStore.AddConnection(&connectionstore.Connection{
				Ch:                      ch,
				LabelAcceptanceCriteria: incomingSubscribe.Criteria,
			})

			uuids = append(uuids, uuid)
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

			connectionStore.SendMessage(incomingMessage.Message)
		default:
			log.Println("unknown type")
			continue
		}
	}

	fmt.Println("exiting websocket connection")
}

func setCors(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
