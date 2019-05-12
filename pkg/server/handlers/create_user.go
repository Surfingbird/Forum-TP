package handlers

import (
	"net/http"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	nickname := c.Param("nickname")

	user := api.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}
	user.Nickname = nickname

	status := models.CreateUser(&user)
	if status == http.StatusConflict {
		users := models.SelectConflictUsers(user.Nickname, user.Email)
		c.JSON(http.StatusConflict, users)

		return
	}

	c.JSON(http.StatusCreated, user)
}
