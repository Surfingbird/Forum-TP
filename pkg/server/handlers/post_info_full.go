package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func PostFullHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	per := vars["id"]
	id, err := strconv.Atoi(per)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	post, status := models.SelectPost(id)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)

		message := "there is no post with this id"
		error := api.Error{
			Message: message,
		}
		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("PostFullHandler, write json: ", err.Error())
		}

		return
	}

	postFull := map[string]interface{}{}
	postFull["post"] = post

	str := r.URL.RawQuery

	if ok, _ := regexp.Match("user", []byte(str)); ok {
		author, err := models.SelectUser(post.Author)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		postFull["author"] = author
	}
	if ok, _ := regexp.Match("forum", []byte(str)); ok {
		forum, status := models.SelectForum(post.Forum)
		if status == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		postFull["forum"] = forum
	}
	if ok, _ := regexp.Match("thread", []byte(str)); ok {
		thread, status := models.SelectThreadBySlugOrID(strconv.Itoa(int(post.Thread)))
		if status == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		postFull["thread"] = thread
	}

	err = json.NewEncoder(w).Encode(postFull)
	if err != nil {
		log.Fatalln("PostFullHandler, write json: ", err.Error())
	}
}
