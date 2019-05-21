package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"DB_Project_TP/api"
	"DB_Project_TP/config"
)

type SelectThreadParams struct {
	Limit int    `schema:"limit"`
	Since string `schema:"since"`
	Desc  bool   `schema:"desc"`
	Query string `schema:"-"`
}

func CreateThread(t *api.Thread) (status int, lastid int64) {
	if ok := CheckUser(t.Author); !ok {
		status = http.StatusNotFound

		return
	}

	ok, forumSlug := CheckForum(t.Forum)
	if !ok {
		status = http.StatusNotFound

		return
	}
	t.Forum = forumSlug

	if t.Slug != "" {
		if ok := CheckThreadBySlug(t.Slug); ok {
			status = http.StatusConflict

			return
		}
	}

	if t.Created.IsZero() {
		err := config.DB.QueryRow(sqlInsertThread,
			t.Author,
			t.Forum,
			t.Message,
			t.Slug,
			t.Title).Scan(&lastid)
		if err != nil {
			config.Logger.Info(fmt.Sprintf("CreateThread: %v", err.Error()))
			status = http.StatusConflict

			return
		}
	} else {
		err := config.DB.QueryRow(sqlInsertThreadWithTime,
			t.Author,
			t.Created,
			t.Forum,
			t.Message,
			t.Slug,
			t.Title).Scan(&lastid)
		if err != nil {
			config.Logger.Info(fmt.Sprintf("CreateThread", err.Error()))
			status = http.StatusConflict

			return
		}
	}
	status = http.StatusCreated

	return
}

func SelectThreadsByForum(forum string, params SelectThreadParams) (threads []api.Thread, status int) {
	if ok, _ := CheckForum(forum); !ok {
		return threads, http.StatusNotFound
	}

	query := SelectQueryFromParams(params)
	rows, err := config.DB.Query(query, forum)
	if err != nil {
		config.Logger.Fatal(fmt.Sprintf("SelectThreadsByForum (%v): %v\n", query, err.Error()))

		return
	}
	defer rows.Close()

	threads = []api.Thread{}
	for rows.Next() {
		thread := api.Thread{}

		if err := rows.Scan(&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Id,
			&thread.Message,
			&thread.Slug,
			&thread.Title); err != nil {
			config.Logger.Fatal(fmt.Sprintf("SelectThreadsByForum (%v): %v\n", query, err.Error()))
		}

		threads = append(threads, thread)
	}

	return threads, http.StatusOK
}

func SelectThread(title, slug string) (thread api.Thread, status int) {
	row, err := config.DB.Query(sqlSelectThread, title, slug)
	if err != nil {
		config.Logger.Fatal(fmt.Sprintf("SelectThread", err.Error()))

		return
	}
	defer row.Close()

	if !row.Next() {
		return thread, http.StatusNotFound
	}

	err = row.Scan(&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Id,
		&thread.Message,
		&thread.Slug,
		&thread.Title)
	if err != nil {
		config.Logger.Fatalf("SelectThread: %v\n", err.Error())

		return
	}

	return thread, http.StatusOK
}

func SelectQueryFromParams(params SelectThreadParams) string {
	sortType := "ASC"
	Limit := "ALL"
	compare := " >= "
	var Since string

	if params.Desc {
		sortType = "DESC"
		compare = " <= "
	}

	if params.Limit != 0 {
		Limit = strconv.Itoa(params.Limit)
	}

	if params.Since != "" {
		Since = params.Since
	}

	query := `select author, created, forum, id, message, slug, title from threads where forum = $1 `

	if Since != "" {
		query = query + "and created" + compare + "'" + Since + "'" + " "
	}

	query = query + "order by created " + sortType + " " +
		"limit" + " " + Limit

	return query
}

func CheckThreadBySlug(slug string) bool {
	res, err := config.DB.Exec(sqlCheckThreadBySlug, slug)
	if err != nil {
		config.Logger.Fatal("CheckThreadBySlug", err.Error())
	}

	count := res.RowsAffected()

	if count > 1 {
		log.Fatalln("Not uniq thread slug")
	}
	if count == 1 {
		return true
	}

	return false
}

func SelectThreadByTitle(title string) (thread api.Thread, status int) {
	row, err := config.DB.Query(sqlSelectThreadByTitle, title)
	if err != nil {
		config.Logger.Fatal("SelectThread", err.Error())
	}
	defer row.Close()

	if !row.Next() {
		return thread, http.StatusNotFound
	}

	err = row.Scan(&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Id,
		&thread.Message,
		&thread.Slug,
		&thread.Title)
	if err != nil {
		config.Logger.Fatalf("SelectThread: %v\n", err.Error())
	}

	return thread, http.StatusOK
}

func ThreadById(id int64) (thread api.Thread, err error) {
	row, err := config.DB.Query(sqlSelectThreadById, id)
	if err != nil {
		config.Logger.Fatal("SelectThread", err.Error())
	}
	defer row.Close()

	if !row.Next() {
		err = fmt.Errorf("ThreadById: There is no this thread!")
		return
	}

	err = row.Scan(&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Id,
		&thread.Message,
		&thread.Slug,
		&thread.Title,
		&thread.Votes)
	if err != nil {
		config.Logger.Fatalf("SelectThread: %v\n", err.Error())
	}

	return
}

func SelectThreadByTitleAndForum(title string, forum string) (thread api.Thread, status int) {
	row, err := config.DB.Query(sqlSelectThreadByTitleAndForum, title, forum)
	if err != nil {
		config.Logger.Fatal("SelectThread", err.Error())
	}
	defer row.Close()

	if !row.Next() {
		return thread, http.StatusNotFound
	}

	err = row.Scan(&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Id,
		&thread.Message,
		&thread.Slug,
		&thread.Title)
	if err != nil {
		config.Logger.Fatalf("SelectThread: %v\n", err.Error())
	}

	return thread, http.StatusOK
}

func ThreadIdbySlug(slug string) (id, status int) {
	row := config.DB.QueryRow(sqlSelectThreadIdbySlug, slug)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return id, http.StatusNotFound
	}

	return id, http.StatusOK
}

func CheckThreadById(id int) bool {
	per := 0
	row := config.DB.QueryRow(sqlCheckThreadIdbyId, id)
	err := row.Scan(&per)
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func GetForumByThread(threadID uint) (slug string, err error) {
	row := config.DB.QueryRow(sqlGetForumByThread, threadID)
	err = row.Scan(&slug)
	if err == sql.ErrNoRows {
		return slug, errors.New("wrong thread id. Can not search forum!")
	}

	return slug, nil
}

func getThreadID(slugOrID string) (id, status int) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		id, status = ThreadIdbySlug(slugOrID)
		if status == http.StatusNotFound {
			return id, http.StatusNotFound
		}
	}
	if ok := CheckThreadById(id); !ok {
		return id, http.StatusNotFound
	}

	return id, http.StatusOK
}

func SelectThreadBySlugOrID(slugOrID string) (thread api.Thread, status int) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		row := config.DB.QueryRow(sqlSelectThreadBySlug, slugOrID)
		err := row.Scan(&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Id,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
		if err != nil {
			status = http.StatusNotFound

			return
		}

		status = http.StatusOK

		return
	} else {
		row := config.DB.QueryRow(sqlSelectThreadById, id)
		err = row.Scan(&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Id,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes)
		if err != nil {
			status = http.StatusNotFound

			return
		}
	}

	status = http.StatusOK

	return
}

func UpdateThread(updateThread api.ThreadUpdate, slugOrID string) (status int) {
	id, status := getThreadID(slugOrID)
	if status == http.StatusNotFound {
		return http.StatusNotFound
	}

	res, err := config.DB.Exec(sqlUpdateThread, updateThread.Message,
		updateThread.Title, id)
	if err != nil {
		config.Logger.Fatal("UpdateThread: ", err.Error())
	}

	rows := res.RowsAffected()
	if rows != 1 {
		config.Logger.Fatalf("UpdateThread, Invalid update: expected %v, have %v\n", 1, rows)
	}

	return http.StatusOK
}

func ThreadIDFromUrl(slugOrID string) (id int, status int) {
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		id, status = ThreadIdbySlug(slugOrID)
		return id, status
	}

	if ok := CheckThreadById(id); !ok {
		return 0, http.StatusNotFound
	}

	return id, http.StatusOK
}

var sqlInsertThreadWithTime2 = `insert into threads (author, created, forum, message, slug, title)
    values ($1,
      $2,
      $3,
      $4,
      $5,
	  $6)
	  returning id`

//toDO костыль pq: CASE types text and timestamp with time zone cannot be matched
var sqlInsertThreadWithTime = `insert into threads (author, created, forum, message, slug, title)
    values ($1,
      $2,
      $3,
      $4,
      $5,
	  $6)
	  returning id`

//toDO костыль pq: CASE types text and timestamp with time zone cannot be matched
var sqlInsertThread = `insert into threads (author, forum, message, slug, title)
    values ($1,
	  $2,
      $3,
      $4,
	  $5)
	  returning id`

var sqlSelectThread = `select author, created, forum, id, message, slug, title
 from threads where title = $1 or slug = $2`

var sqlCheckThreadBySlug = `select author, created, forum, id, message, slug, title
 from threads where slug = $1`

var sqlSelectThreadByTitle = `select author, created, forum, id, message, slug, title
 from threads where title = $1`

var sqlSelectThreadByTitleAndForum = `select author, created, forum, id, message, slug, title
 from threads where title = $1 and forum = $2`

var sqlSelectThreadById = `select author, created, forum, id, message, slug, title, votes
 from threads where id = $1`

var sqlSelectThreadBySlug = `select author, created, forum, id, message, slug, title, votes
from threads where slug = $1`

var sqlSelectThreadIdbySlug = `select id from threads where slug = $1`

var sqlCheckThreadIdbyId = `select id from threads where id = $1`

var sqlGetForumByThread = `select f.slug from forums f where f.slug = 
(select t.forum from threads t where t.id = $1)`

var sqlUpdateThread = `update threads  set message = (case
            when $1 = '' then message
             else $1 end),
    title = (case
            when $2 = '' then title
            else $2 end)
where id = $3`
