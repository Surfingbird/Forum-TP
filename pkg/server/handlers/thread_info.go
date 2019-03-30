package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func ThreadInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	treadID := vars["slug_or_id"]

	thread, status := models.SelectThreadBySlugOrID(treadID)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		message := "We can not find this thread!"
		error := api.Error{
			Message: message,
		}

		err := json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("ThreadInfo, write json: ", err.Error())
		}

		return
	}

	err := json.NewEncoder(w).Encode(thread)
	if err != nil {
		log.Fatalln("ThreadInfo, write json: ", err.Error())
	}
}
