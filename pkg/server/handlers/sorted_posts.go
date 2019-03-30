package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func SortedPostsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrID := vars["slug_or_id"]
	threadID, status := models.ThreadIDFromUrl(slugOrID)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		message := "There is not this thread"
		error := api.Error{
			Message: message,
		}
		err := json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("SortedPostsHandler, write json: ", err.Error())
		}

		return
	}

	params := &api.PostsSorted{}
	err := decoder.Decode(params, r.URL.Query())
	if err != nil {
		log.Fatalln("SortedPostsHandler", err.Error())
	}

	posts := models.SortedPosts(params, threadID)

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Fatalln("SortedPostsHandler, write json: ", err.Error())
	}
}
