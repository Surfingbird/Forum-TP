package handlers

import (
	"log"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func SortedPostsHandler(c *gin.Context) {
	slugOrID := c.Param("slug_or_id")
	threadID, status := models.ThreadIDFromUrl(slugOrID)
	if status == http.StatusNotFound {
		message := "There is not this thread"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	params := &api.PostsSorted{}
	err := decoder.Decode(params, c.Request.URL.Query())
	if err != nil {
		log.Fatalln("SortedPostsHandler", err.Error())
	}

	posts := models.SortedPosts(params, threadID)

	c.JSON(http.StatusOK, posts)
}
