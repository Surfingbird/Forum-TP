package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"DB_Project_TP/api"
	"DB_Project_TP/pkg/server/models"

	"github.com/gin-gonic/gin"
)

func PostFullHandler(c *gin.Context) {
	per := c.Param("id")

	id, err := strconv.Atoi(per)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	post, status := models.SelectPost(id)
	if status == http.StatusNotFound {
		message := "there is no post with this id"
		error := api.Error{
			Message: message,
		}

		c.JSON(http.StatusNotFound, error)

		return
	}

	postFull := map[string]interface{}{}
	postFull["post"] = post

	str := c.Request.URL.RawQuery

	if ok, _ := regexp.Match("user", []byte(str)); ok {
		author, err := models.SelectUser(post.Author)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		postFull["author"] = author
	}
	if ok, _ := regexp.Match("forum", []byte(str)); ok {
		forum, status := models.SelectForum(post.Forum)
		if status == http.StatusNotFound {
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		postFull["forum"] = forum
	}
	if ok, _ := regexp.Match("thread", []byte(str)); ok {
		thread, status := models.SelectThreadBySlugOrID(strconv.Itoa(int(post.Thread)))
		if status == http.StatusNotFound {
			c.AbortWithStatus(http.StatusNotFound)

			return
		}

		postFull["thread"] = thread
	}

	c.JSON(http.StatusOK, postFull)
}
