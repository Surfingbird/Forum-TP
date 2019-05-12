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

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	threadId, status := models.ThreadIDFromUrl(slugOrId)
	if status == http.StatusNotFound {
		message := "There is no this thread"
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("CreatePostHandler, read json: ", err.Error())
	}

	posts := []api.Post{}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Fatalln("CreatePostHandler, unmarshal json: ", err.Error())
	}

	status, postsId := models.CreatePost(posts, threadId)
	if status == http.StatusConflict {
		message := "There is no post's parent"
		error := api.Error{
			Message: message,
		}

		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("ForumsBranchsHandler, write json: ", err.Error())
		}

		return
	}

	if status == http.StatusNotFound {
		message := "There is no this branch"
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

	postsFull := models.SelectCreatedPosts(postsId)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(postsFull)
	if err != nil {
		log.Fatalln("ForumsBranchsHandler, write json: ", err.Error())
	}
}
