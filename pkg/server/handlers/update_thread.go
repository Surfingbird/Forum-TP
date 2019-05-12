package handlers

import (
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func UpdateBranchHandler(c *gin.Context) {
	treadID := c.Param("slug_or_id")

	updateThread := api.ThreadUpdate{}
	err := c.BindJSON(&updateThread)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	status := models.UpdateThread(updateThread, treadID)
	if status == http.StatusNotFound {
		message := "We can not finc this thread"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	thread, _ := models.SelectThreadBySlugOrID(treadID)

	c.JSON(http.StatusOK, thread)
}
