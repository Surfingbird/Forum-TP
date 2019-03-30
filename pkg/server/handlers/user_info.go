package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	user, err := models.SelectUser(nickname)
	if err != nil {
		message := fmt.Sprintf("Can't find user with nickname %v", nickname)
		error := api.Error{
			Message: message,
		}

		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("CreateUserHandler, write json: ", err.Error())
		}

		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Fatalln("CreateUserHandler, write json: ", err.Error())
	}
}
