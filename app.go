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

type App struct {
	Router *mux.Router
}

type EmailJson struct {
	To string `json:"To"`
	From string `json:"From"`
	Date string `json:"Date"`
	Subject string `json:"Subject"`
	MessageID string `json:"Message-ID"`
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/parse", ParseEmailEndpoint).Methods("POST")
}

func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func ParseEmailEndpoint(w http.ResponseWriter, req *http.Request) {
	//params  := mux.Vars(req)
	body, err := ioutil.ReadAll(req.Body)
	msg, err := mail.ReadMessage(bytes.NewBuffer([]byte(body)))
	var emailJson EmailJson
	if err != nil {
		fmt.Println("Invalid Email Format")
		respondWithError(w, http.StatusBadRequest, "Invalid Email Format")
		return
	}
	fmt.Println(msg.Header.Get("Message-ID"))
	emailJson.To = msg.Header.Get("To")
	emailJson.From = msg.Header.Get("From")
	emailJson.Date = msg.Header.Get("Date")
	emailJson.Subject = msg.Header.Get("Subject")
	emailJson.MessageID = msg.Header.Get("Message-ID")

	respondWithJSON(w, http.StatusOK, emailJson)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
