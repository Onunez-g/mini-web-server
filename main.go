package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type message struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

var messages []message = []message{
	{
		Id:      1,
		Message: "Hola",
	},
	{
		Id:      2,
		Message: "Mundo",
	},
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/messages", GetMessages).Methods("GET")
	r.HandleFunc("/messages", createMessage).Methods("POST")
	r.HandleFunc("/messages/{id}", updateMessage).Methods("PUT")
	r.HandleFunc("/messages/{id}", deleteMessage).Methods("DELETE")

	apiPort := ":8000"
	fmt.Printf("Listening server at port %s\n", apiPort)
	err := http.ListenAndServe(apiPort, r)
	if err != nil {
		log.Panicln(err.Error())
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(&messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(fmt.Errorf("An error has ocurred: %s", err.Error()))
	}
	if _, err := w.Write([]byte(response)); err != nil {
		log.Fatal(err.Error())
	}
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var msg message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(fmt.Errorf("An error has ocurred: %s", err.Error()))
		log.Fatal(w.Write([]byte("Wrong format")))
		return
	}
	for _, v := range messages {
		if v.Id == msg.Id {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatal(fmt.Errorf("An error has ocurred: %s", err.Error()))
			log.Fatal(w.Write([]byte("Existing Id")))
			return
		}
	}
	messages = append(messages, msg)
	response, err := json.Marshal(&msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(fmt.Errorf("An error has ocurred: %s", err.Error()))
		return
	}
	if _, err := w.Write([]byte(response)); err != nil {
		log.Fatal(err.Error())
	}
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var msg message = message{}
	var updatedMsg message
	err := json.NewDecoder(r.Body).Decode(&updatedMsg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(fmt.Errorf("An error has ocurred: %s", err.Error()))
		log.Fatal(w.Write([]byte("Wrong format")))
		return
	}

	for k, v := range messages {
		if strconv.Itoa(v.Id) == params["id"] {
			updatedMsg.Id = v.Id
			messages[k] = updatedMsg
			msg = updatedMsg
		}
	}
	if msg.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w.Write([]byte("Message not Found")))
		return
	}
	response, err := json.Marshal(&msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(fmt.Errorf("An error has ocurred: %s", err.Error()))
		return
	}
	if _, err := w.Write([]byte(response)); err != nil {
		log.Fatal(err.Error())
	}

}

func deleteMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var msg message = message{}

	for k, v := range messages {
		if strconv.Itoa(v.Id) == params["id"] {
			messages = append(messages[:k], messages[k+1:]...)
			msg = v
		}
	}
	if msg.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w.Write([]byte("Message not Found")))
		return
	}
	response, err := json.Marshal(&msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(fmt.Errorf("An error has ocurred: %s", err.Error()))
		return
	}
	if _, err := w.Write([]byte(response)); err != nil {
		log.Fatal(err.Error())
	}
}
