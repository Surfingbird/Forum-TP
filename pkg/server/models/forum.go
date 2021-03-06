package models

import (
	"DB_Project_TP/api"
	"DB_Project_TP/config"
	"net/http"
)

func CreateForum(f *api.Forum) (status int) {
	if ok := CheckUser(f.User); !ok {
		return http.StatusNotFound
	}

	regNick, _ := RegNickname(f.User)
	f.User = regNick

	_, err := config.DB.Exec(sqlInsertForum,
		f.Slug,
		f.Threads,
		f.Title,
		f.User)
	if err != nil {
		return http.StatusConflict
	}

	return http.StatusCreated
}

func SelectForum(slug string) (forum api.Forum, status int) {
	err := config.DB.QueryRow(sqlSelectForum, slug).Scan(
		&forum.Slug,
		&forum.Threads,
		&forum.Title,
		&forum.Posts,
		&forum.User)

	if err != nil {
		status = http.StatusNotFound

		return
	}

	status = http.StatusOK

	return
}

func CheckForum(slug string) (ok bool, forumSlug string) {
	row := config.DB.QueryRow(sqlCheckForum, slug)
	err := row.Scan(&forumSlug)
	if err != nil {
		return false, forumSlug
	}

	return true, forumSlug
}

var sqlInsertForum = `insert into forums (slug, threads, title, user_f)
values ($1, $2, $3, $4)`

var sqlSelectForum = `select slug,
  (select count(*) from threads where forum = $1) as threads,
  title,
  (select count(*) from posts where forum = $1) as post,
  user_f
from forums where slug = $1`

var sqlCheckForum = `select slug from forums where slug = $1`
