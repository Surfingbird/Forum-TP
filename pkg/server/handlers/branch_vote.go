package handlers

import (
	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BranchVoteHandler(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")
	treadID, status := models.ThreadIDFromUrl(slugOrId)
	if status == http.StatusNotFound {
		message := "Can not find thread with this id or slug!"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	vote := api.Vote{}
	err := c.BindJSON(&vote)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	status = models.VoteBranch(vote, treadID)
	if status == http.StatusNotFound {
		message := "Can not find thread with this id!"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	thread, _ := models.ThreadById(int64(treadID))

	c.JSON(http.StatusOK, thread)
}
