package Bithose

import (
	"encoding/json"
	"github.com/JonathanRosado/Bithose/connectionstore"
	"net/http"
)

func StatsHandler(writer http.ResponseWriter, request *http.Request) {
	setCors(writer)

	connectionStore := connectionstore.GetMapStore()

	jsonData, err := json.Marshal(connectionStore.Stats())
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonData)
}
