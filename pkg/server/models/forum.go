package models

import (
	"DB_Project_TP/api"
	"DB_Project_TP/config"
	"database/sql"
	"log"
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
	row, err := config.DB.Query(sqlSelectForum, slug)
	if err != nil {
		log.Fatalln("SelectForum", err.Error())
	}
	defer row.Close()

	if !row.Next() {
		return forum, http.StatusNotFound
	}

	if err := row.Scan(
		&forum.Slug,
		&forum.Threads,
		&forum.Title,
		&forum.Posts,
		&forum.User); err != nil {
		log.Fatalf("SelectUser: %v\n", err.Error())
	}

	return forum, http.StatusOK
}

func CheckForum(slug string) (ok bool, forumSlug string) {
	row := config.DB.QueryRow(sqlCheckForum, slug)
	err := row.Scan(&forumSlug)
	if err == sql.ErrNoRows {
		return false, forumSlug
	}

	return true, forumSlug
}

var sqlInsertForum = `insert into project_bd.forums (slug, threads, title, user_f)
values ($1, $2, $3, $4)`

var sqlSelectForum = `select slug,
  (select count(*) from project_bd.threads where forum = $1) as threads,
  title,
  (select count(*) from project_bd.posts where forum = $1) as post,
  user_f
from project_bd.forums where slug = $1`

var sqlCheckForum = `select slug from project_bd.forums where slug = $1`
