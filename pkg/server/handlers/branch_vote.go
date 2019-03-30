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

func BranchVoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	treadID := vars["slug_or_id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("CreatePostHandler, read json: ", err.Error())
	}

	vote := api.Vote{}
	err = json.Unmarshal(body, &vote)
	if err != nil {
		log.Fatalln("CreatePostHandler, unmarshal json: ", err.Error())
	}

	status := models.VoteBranch(vote, treadID)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		message := "Can not find thread with this id!"
		error := api.Error{
			Message: message,
		}

		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("BranchVoteHandler, write json: ", err.Error())
		}

		return
	}

	thread, _ := models.SelectThreadBySlugOrID(treadID)
	err = json.NewEncoder(w).Encode(thread)
	if err != nil {
		log.Fatalln("BranchVoteHandler, write json: ", err.Error())
	}
}
