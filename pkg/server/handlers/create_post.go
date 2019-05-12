package handlers

import (
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")
	threadId, status := models.ThreadIDFromUrl(slugOrId)
	if status == http.StatusNotFound {
		message := "There is no this thread"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	posts := []api.Post{}
	err := c.BindJSON(&posts)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	status, postsId := models.CreatePost(posts, threadId)
	if status == http.StatusConflict {
		message := "There is no post's parent"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusConflict, error)

		return
	}

	if status == http.StatusNotFound {
		message := "There is no this branch"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	postsFull := models.SelectCreatedPosts(postsId)

	c.JSON(http.StatusCreated, postsFull)
}
