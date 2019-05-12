package handlers

import (
	"net/http"

	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func DBInfoHandler(c *gin.Context) {
	info := models.GetDBInfo()

	c.JSON(http.StatusOK, info)
}
