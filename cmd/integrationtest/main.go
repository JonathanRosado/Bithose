package main

import (
	"github.com/JonathanRosado/Bithose"
	"github.com/JonathanRosado/Bithose/cmd/integrationtest/stress"
	"log"
	"net/http"
)

func main() {

	go func() {
		http.HandleFunc("/", Bithose.WsHandler)
		http.HandleFunc("/stats", Bithose.StatsHandler)
		log.Fatal(http.ListenAndServe(":80", nil))
	}()

	stress.Run()
}
