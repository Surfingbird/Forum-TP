package handlers

import (
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func ForumsUsersHandlers(c *gin.Context) {
	slug := c.Param("slug")

	params := api.ForumsUsersQuery{}
	err := decoder.Decode(&params, c.Request.URL.Query())
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	posts, status := models.SelectForumsUsers(params, slug)
	if status == http.StatusNotFound {
		message := "there is no this forum"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	c.JSON(http.StatusOK, posts)
}
