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

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user := api.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("CreateUserHandler, read json: ", err.Error())
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatalln("CreateUserHandler, unmarshal json: ", err.Error())
	}
	user.Nickname = nickname

	status := models.CreateUser(&user)
	if status == http.StatusCreated {
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Fatalln("CreateUserHandler, write json: ", err.Error())
		}

		return
	}

	if status == http.StatusConflict {
		w.WriteHeader(http.StatusConflict)

		users := models.SelectConflictUsers(user.Nickname, user.Email)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			log.Fatalln("CreateUserHandler, write json: ", err.Error())
		}

		return
	}

	log.Fatalln("CreateUserHandler: 500, internal error")
}
