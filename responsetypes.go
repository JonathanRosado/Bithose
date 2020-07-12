package Bithose

type SendMessageResponse struct {
	NumberOfSents    int    `json:"number_of_sents"`
	NumberOfTimeouts int    `json:"number_of_timeouts"`
	Error            string `json:"error"`
}

type SubscribeResponse struct {
	Uuid  string `json:"uuid"`
	Error string `json:"error"`
}
