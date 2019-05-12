package handlers

import (
	"fmt"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/config"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func CreateForumHandler(c *gin.Context) {
	forum := api.Forum{}
	err := c.BindJSON(&forum)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	status := models.CreateForum(&forum)
	if status == http.StatusNotFound {
		message := fmt.Sprint("This user can not create forum!")
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	if status == http.StatusConflict {
		forum, status := models.SelectForum(forum.Slug)
		if status == http.StatusNotFound {
			config.Logger.Fatal("Can not find already exists forum")
		}

		c.JSON(http.StatusConflict, forum)

		return
	}

	c.JSON(http.StatusCreated, forum)
}
