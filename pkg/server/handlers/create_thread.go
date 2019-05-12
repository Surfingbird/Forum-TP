package handlers

import (
	"fmt"
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func CreateThreadHandler(c *gin.Context, slug string) {
	thread := api.Thread{}
	err := c.BindJSON(&thread)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}
	thread.Forum = slug

	status, id := models.CreateThread(&thread)
	if status == http.StatusNotFound {
		message := fmt.Sprint("Can not find author or forum!")
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	// TODO исправить узкое место
	if status == http.StatusConflict {
		thread, _ = models.SelectThread(thread.Title, thread.Slug)
		c.JSON(http.StatusConflict, thread)

		return
	}

	thread, err = models.ThreadById(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		log.Fatalln("Can not search created thread")

		return
	}

	c.JSON(http.StatusCreated, thread)
}
