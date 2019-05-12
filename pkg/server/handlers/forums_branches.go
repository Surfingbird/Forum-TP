package handlers

import (
	"log"
	"net/http"
	"strconv"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func ForumsBranchsHandler(c *gin.Context) {
	slug := c.Param("slug")

	err := c.Request.ParseForm()
	if err != nil {
		log.Fatal("ParseForm", err.Error())
	}

	since := c.Request.FormValue("since")
	params := models.SelectThreadParams{
		Since: since,
	}

	rowLimit, err := strconv.Atoi(c.Request.FormValue("limit"))
	if err == nil {
		params.Limit = rowLimit
	}

	rowDesc, err := strconv.ParseBool(c.Request.FormValue("desc"))
	if err == nil {
		params.Desc = rowDesc
	}

	threads, status := models.SelectThreadsByForum(slug, params)
	if status == http.StatusNotFound {
		message := "Can't find threads with this slug"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	c.JSON(http.StatusOK, threads)
}
