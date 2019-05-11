package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	thread := api.Thread{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("CreateThreadHandler, read json: ", err.Error())
	}

	err = json.Unmarshal(body, &thread)
	if err != nil {
		log.Fatalln("CreateThreadHandler, unmarshal json: ", err.Error())
	}
	thread.Forum = slug

	status, id := models.CreateThread(&thread)
	if status == http.StatusNotFound {
		message := fmt.Sprint("Can not find author or forum!")
		error := api.Error{
			Message: message,
		}

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("CreateThreadHandler, write json: ", err.Error())
		}

		return
	}

	if status == http.StatusConflict {
		w.WriteHeader(http.StatusConflict)

		thread, _ = models.SelectThread(thread.Title, thread.Slug)
		err = json.NewEncoder(w).Encode(thread)
		if err != nil {
			log.Fatalln("CreateThreadHandler, write json: ", err.Error())
		}

		return
	}

	// bug, select thread by forum and title !!!! here
	thread, err = models.ThreadById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Fatalln("Can not search created thread")

		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(thread)
	if err != nil {
		log.Fatalln("CreateThreadHandler, write json: ", err.Error())
	}
}
