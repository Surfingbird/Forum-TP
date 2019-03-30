package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func ForumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	forum, status := models.SelectForum(slug)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)

		message := fmt.Sprint("There is no forum with this slug")
		error := api.Error{
			Message: message,
		}

		err := json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("ForumHandler, write json: ", err.Error())
		}

		return
	}

	err := json.NewEncoder(w).Encode(forum)
	if err != nil {
		log.Fatalln("ForumHandler, write json: ", err.Error())
	}
}
