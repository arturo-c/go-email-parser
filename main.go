package main

import (
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"fmt"
	"net/mail"
	"github.com/gorilla/mux"
	"io/ioutil"
)

type EmailJson struct {
	To string `json:"To"`
	From string `json:"From"`
	Date string `json:"Date"`
	Subject string `json:"Subject"`
	MessageID string `json:"Message-ID"`
}

func ParseEmailEndpoint(w http.ResponseWriter, req *http.Request) {
	//params  := mux.Vars(req)
	body, err := ioutil.ReadAll(req.Body)
	msg, err := mail.ReadMessage(bytes.NewBuffer([]byte(body)))
	var emailJson EmailJson
	if err != nil {
		panic(err)
	}
	fmt.Println(msg.Header.Get("Message-ID"))
	emailJson.To = msg.Header.Get("To")
	emailJson.From = msg.Header.Get("From")
	emailJson.Date = msg.Header.Get("Date")
	emailJson.Subject = msg.Header.Get("Subject")
	emailJson.MessageID = msg.Header.Get("Message-ID")

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(emailJson)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/parse", ParseEmailEndpoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":12345", router))
}
