package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func ForumsUsersHandlers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	params := api.ForumsUsersQuery{}
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		log.Fatalln("ForumsUsersHandlers: Can not parse url")
	}

	posts, status := models.SelectForumsUsers(params, slug)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		message := "there is no this forum"
		error := api.Error{
			Message: message,
		}

		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("ForumsUsersHandlers, write json: ", err.Error())
		}

		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Fatalln("ForumsUsersHandlers, write json: ", err.Error())
	}
}
