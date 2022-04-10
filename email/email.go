package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Sender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type To struct {
	Email string `json:"email"`
}

type MailReq struct {
	Sender      Sender `json:"sender"`
	To          []To   `json:"to"`
	Subject     string `json:"subject"`
	HtmlContent string `json:"htmlContent"`
}

func SendMail(emailTo string, link string) error {
	c := http.Client{Timeout: time.Duration(0) * time.Second}

	endpoint := os.Getenv("SEND_IN_BLUE_ENDPOINT")

	API_KEY := os.Getenv("SEND_IN_BLUE_API_KEY")

	emailFrom := os.Getenv("EMAIL_FROM")

	emailFromName := os.Getenv("EMAIL_FROM_NAME")

	subject := "Email Verification from Golang Auth Service"

	information := "Click this link to verify your account[Golang Auth]: " + link

	m := MailReq{
		Sender: Sender{
			Email: emailFrom,
			Name:  emailFromName,
		},
		To: []To{
			To{
				Email: emailTo,
			},
		},
		Subject:     subject,
		HtmlContent: information,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	req, err := http.NewRequest("POST", endpoint, b)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-key", API_KEY)

	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	return err
}
