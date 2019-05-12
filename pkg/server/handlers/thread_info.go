package handlers

import (
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func ThreadInfo(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	thread, status := models.SelectThreadBySlugOrID(slugOrId)
	if status == http.StatusNotFound {
		message := "We can not find this thread!"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	c.JSON(http.StatusOK, thread)
}
