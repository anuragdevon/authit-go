package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InternalRequest(payload map[string]interface{}, reqtype string, endpoint string) (*http.Response, error) {

	json_payload, err := json.Marshal(payload)
	if err != nil {
		log.Panic(err.Error())
	}
	var resp *http.Response

	if reqtype == "POST" {
		resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(json_payload))

	} else {
		resp, err = http.Post(endpoint, "application/json", bytes.NewBuffer(json_payload))
	}
	return resp, err

}

func Response(c *gin.Context, code int, obj interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(code, obj)
}