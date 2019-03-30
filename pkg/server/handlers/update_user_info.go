package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gorilla/mux"
)

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	update := api.UpdateUser{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("UpdateProfileHandler, read json: ", err.Error())
	}

	err = json.Unmarshal(body, &update)
	if err != nil {
		log.Fatalln("UpdateProfileHandler, unmarshal json: ", err.Error())
	}

	user, status := models.UpdateUser(&update, nickname)
	if status == http.StatusNotFound {
		message := fmt.Sprintf("Can't find user with nickname %v", nickname)
		error := api.Error{
			Message: message,
		}

		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("UpdateProfileHandler, write json: ", err.Error())
		}

		return
	}

	if status == http.StatusConflict {
		message := fmt.Sprintf("Can't updaet %v with this data", nickname)
		error := api.Error{
			Message: message,
		}

		w.WriteHeader(http.StatusConflict)
		err = json.NewEncoder(w).Encode(error)
		if err != nil {
			log.Fatalln("UpdateProfileHandler, write json: ", err.Error())
		}

		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Fatalln("UpdateProfileHandler, write json: ", err.Error())
	}
}
