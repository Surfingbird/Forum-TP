package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func UpdateBranchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	treadID := vars["slug_or_id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("UpdateBranchHandler, read json: ", err.Error())
	}

	updateThread := api.ThreadUpdate{}
	err = json.Unmarshal(body, &updateThread)
	if err != nil {
		log.Fatalln("UpdateBranchHandler, unmarshal json: ", err.Error())
	}

	status := models.UpdateThread(updateThread, treadID)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		message := "We can not finc this thread"
		error := api.Error{
			Message: message,
		}

		err := json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("UpdateBranchHandler, write json: ", err.Error())
		}

		return
	}

	thread, _ := models.SelectThreadBySlugOrID(treadID)
	err = json.NewEncoder(w).Encode(thread)
	if err != nil {
		log.Fatalln("UpdateBranchHandler, write json: ", err.Error())
	}
}
