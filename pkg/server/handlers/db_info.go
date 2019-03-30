package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"DB_Project_TP/pkg/server/models"
)

func DBInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := models.GetDBInfo()
	err := json.NewEncoder(w).Encode(info)
	if err != nil {
		log.Fatalln("DBInfoHandler, write json: ", err.Error())
	}
}
