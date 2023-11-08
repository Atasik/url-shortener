package v1

import (
	"encoding/json"
	"log"
	"net/http"
)

type errResponse struct {
	Message string `json:"message"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type linkResponse struct {
	Link string `json:"link"`
}

func newErrorResponse(w http.ResponseWriter, msg string, status int) {
	resp, err := json.Marshal(errResponse{msg})
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}
