package handler

import (
	"net/http"
)

func PingHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	pingResponse := map[string]interface{}{
		"message": "PONG!",
	}
	sendResponseToClient(responseWriter, "", pingResponse, nil, 200)
}
