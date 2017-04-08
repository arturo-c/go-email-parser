package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"time"
)

type App struct {
	Router      *mux.Router
	RedisClient *redis.Client
}

type EmailJson struct {
	To        string `json:"To"`
	From      string `json:"From"`
	Date      string `json:"Date"`
	Subject   string `json:"Subject"`
	MessageID string `json:"Message-ID"`
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/parse", a.ParseEmailEndpoint).Methods("POST")
}

func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, a.Router))
}

func (a *App) Initialize() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	a.RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) ParseEmailEndpoint(w http.ResponseWriter, req *http.Request) {
	//params  := mux.Vars(req)
	cache, err := a.RedisClient.Get("some-cache").Result()
	fmt.Println(cache)
	if err == redis.Nil {
		time.Sleep(time.Second * 2)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println("Invalid Email Format")
			respondWithError(w, http.StatusBadRequest, "Invalid Email Format")
			return
		}
		msg, err := mail.ReadMessage(bytes.NewBuffer([]byte(body)))
		var emailJson EmailJson
		if err != nil {
			log.Warn("Invalid Email Format")
			respondWithError(w, http.StatusBadRequest, "Invalid Email Format")
			return
		}

		emailJson.To = msg.Header.Get("To")
		emailJson.From = msg.Header.Get("From")
		emailJson.Date = msg.Header.Get("Date")
		emailJson.Subject = msg.Header.Get("Subject")
		emailJson.MessageID = msg.Header.Get("Message-ID")

		json_response, _ := json.Marshal(emailJson)
		err = a.RedisClient.Set("some-cache", json_response, 0).Err()
		if err != nil {
			panic(err)
		}

		log.WithFields(log.Fields{"email-message": emailJson}).Info("New Email Parse")

		respondWithJSON(w, http.StatusOK, emailJson)
	} else if err != nil {
		panic(err)
	} else {
		var cache_response EmailJson
		json.Unmarshal([]byte(cache), &cache_response)
		respondWithJSON(w, http.StatusOK, cache_response)
	}
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
