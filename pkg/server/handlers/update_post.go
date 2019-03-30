package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	per := vars["id"]
	id, err := strconv.Atoi(per)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("UpdatePostHandler, read json: ", err.Error())
	}

	update := api.PostUpdaet{}
	err = json.Unmarshal(body, &update)
	if err != nil {
		log.Fatalln("UpdatePostHandler, unmarshal json: ", err.Error())
	}

	status := models.UpdatePost(id, update)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		msg := "Can not find post with this ID"
		error := api.Error{
			Message: msg,
		}
		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("UpdatePostHandler, write json: ", err.Error())
		}

		return
	}

	post, _ := models.SelectPost(id)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		log.Fatalln("UpdatePostHandler, write json: ", err.Error())
	}
}
