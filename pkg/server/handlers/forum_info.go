package handlers

import (
	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ForumHandler(c *gin.Context) {
	slug := c.Param("slug")

	forum, status := models.SelectForum(slug)
	if status == http.StatusNotFound {
		message := fmt.Sprint("There is no forum with this slug")
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	c.JSON(http.StatusOK, forum)
}
