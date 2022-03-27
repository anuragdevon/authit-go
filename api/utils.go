package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func InternalRequest(payload map[string]interface{}, reqtype string, endpoint string) (error, *http.Response) {

	json_payload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	var resp *http.Response

	if reqtype == "POST" {
		resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(json_payload))

	} else {
		resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(json_payload))
	}
	return err, resp

}
