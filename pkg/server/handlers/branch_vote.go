package handlers

import (
	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func BranchVoteHandler(c *gin.Context) {
	slugOrId := c.Param("slug_or_id")

	vote := api.Vote{}
	err := c.BindJSON(&vote)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	wg := &sync.WaitGroup{}
	threadChan := make(chan api.Thread, 1)
	okUserChan := make(chan bool, 1)

	wg.Add(1)
	go func(nickname string, okUserChan chan bool, wg *sync.WaitGroup) {
		defer wg.Done()

		ok := models.CheckUser(nickname)
		okUserChan <- ok

	}(vote.Nickname, okUserChan, wg)

	wg.Add(1)
	go func(slugOrId string, threadChan chan api.Thread, wg *sync.WaitGroup) {
		defer wg.Done()

		tread, status := models.SelectThreadBySlugOrID(slugOrId)
		if status == http.StatusNotFound {
			threadChan <- api.Thread{}

			return
		}

		threadChan <- tread
	}(slugOrId, threadChan, wg)

	wg.Wait()
	ok := <-okUserChan
	thread := <-threadChan

	if thread.Id == 0 || ok == false {
		message := "Can not find user or thread"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	sum := models.VoteBranch(vote, thread.Id)

	thread.Votes = int64(sum)

	c.JSON(http.StatusOK, thread)
}
