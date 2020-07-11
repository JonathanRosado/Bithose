package Bithose

import (
	"Bithose/connectionstore"
	"encoding/json"
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
