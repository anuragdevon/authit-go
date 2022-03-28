package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func InternalRequest(payload map[string]interface{}, reqtype string, endpoint string) (*http.Response, error) {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Panic(err.Error())
	}
	var resp *http.Response

	if reqtype == "POST" {
		resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))

	} else {
		resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	}
	return resp, err
}
