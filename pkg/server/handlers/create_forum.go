package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"
)

func CreateForumHandler(w http.ResponseWriter, r *http.Request) {
	forum := api.Forum{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("CreateForumHandler, read json: ", err.Error())
	}

	err = json.Unmarshal(body, &forum)
	if err != nil {
		log.Fatalln("CreateForumHandler, unmarshal json: ", err.Error())
	}

	status := models.CreateForum(&forum)
	if status == http.StatusNotFound {
		message := fmt.Sprint("This user can not create forum!")
		error := api.Error{
			Message: message,
		}

		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("CreateForumHandler, write json: ", err.Error())
		}

		return
	}

	if status == http.StatusConflict {
		forum, status := models.SelectForum(forum.Slug)
		if status == http.StatusNotFound {
			log.Fatalln("Can not find already exists forum")
		}

		w.WriteHeader(http.StatusConflict)
		err = json.NewEncoder(w).Encode(forum)
		if err != nil {
			log.Fatalln("CreateForumHandler, write json: ", err.Error())
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(forum)
	if err != nil {
		log.Fatalln("CreateForumHandler, write json: ", err.Error())
	}
}
