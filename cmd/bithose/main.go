package main

import (
	"flag"
	"github.com/JonathanRosado/Bithose"
	"log"
	"net/http"
)

var (
	hostname string
)

func init() {
	flag.StringVar(&hostname, "hostname", ":9483", "hostname for the bithose server to "+
		"run on [:9483]")
}

func main() {
	http.HandleFunc("/", Bithose.WsHandler)
	http.HandleFunc("/stats", Bithose.StatsHandler)
	log.Fatal(http.ListenAndServe(hostname, nil))
}
