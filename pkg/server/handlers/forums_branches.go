package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func ForumsBranchsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	//toDo вынести этот кусок в gorilla schema
	err := r.ParseForm()
	if err != nil {
		log.Fatalln("ParseForm", err.Error())
	}

	since := r.FormValue("since")
	params := models.SelectThreadParams{
		Since: since,
	}

	rowLimit, err := strconv.Atoi(r.FormValue("limit"))
	if err == nil {
		params.Limit = rowLimit
	}

	rowDesc, err := strconv.ParseBool(r.FormValue("desc"))
	if err == nil {
		params.Desc = rowDesc
	}
	//toDo вынести этот кусок в gorilla schema

	threads, status := models.SelectThreadsByForum(slug, params)
	if status == http.StatusNotFound {
		message := "Can't find threads with this slug"
		error := api.Error{
			Message: message,
		}

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("ForumsBranchsHandler, write json: ", err.Error())
		}

		return
	}

	err = json.NewEncoder(w).Encode(threads)
	if err != nil {
		log.Fatalln("ForumsBranchsHandler, write json: ", err.Error())
	}
}
