package handlers

import (
	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdatePostHandler(c *gin.Context) {
	per := c.Param("id")
	id, err := strconv.Atoi(per)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	update := api.PostUpdaet{}
	err = c.BindJSON(&update)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	status := models.UpdatePost(id, update)
	if status == http.StatusNotFound {
		msg := "Can not find post with this ID"
		error := api.Error{
			Message: msg,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	post, _ := models.SelectPost(id)

	c.JSON(http.StatusOK, post)
}
