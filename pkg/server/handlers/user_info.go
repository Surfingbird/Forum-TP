package handlers

import (
	"fmt"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func ProfileHandler(c *gin.Context) {
	nickname := c.Param("nickname")

	user, err := models.SelectUser(nickname)
	if err != nil {
		message := fmt.Sprintf("Can't find user with nickname %v", nickname)
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	c.JSON(http.StatusOK, user)
}
