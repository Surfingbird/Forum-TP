package handlers

import (
	"fmt"
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func UpdateProfileHandler(c *gin.Context) {
	nickname := c.Param("nickname")

	update := api.UpdateUser{}
	err := c.BindJSON(&update)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	user, status := models.UpdateUser(&update, nickname)
	if status == http.StatusNotFound {
		message := fmt.Sprintf("Can't find user with nickname %v", nickname)
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	if status == http.StatusConflict {
		message := fmt.Sprintf("Can't updaet %v with this data", nickname)
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusConflict, error)

		return
	}

	c.JSON(http.StatusOK, user)
}
