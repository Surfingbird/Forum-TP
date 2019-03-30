package handlers

import (
	"DB_Project_TP/pkg/server/models"
	"net/http"
)

func ClearDB(w http.ResponseWriter, r *http.Request) {
	models.TruncateAllTables()
}
