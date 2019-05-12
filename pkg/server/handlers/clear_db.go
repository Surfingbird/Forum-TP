package handlers

import (
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func ClearDB(c *gin.Context) {
	models.TruncateAllTables()
}
