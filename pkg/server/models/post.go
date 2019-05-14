package models

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"DB_Project_TP/api"
	"DB_Project_TP/config"
)

func CreatePost(posts []api.Post, treadID int) (status int, postsId []int) {
	postsId = []int{}

	tx, _ := config.DB.Begin()
	stmnt, err := tx.Prepare(sqlInsertPost)
	if err != nil {
		tx.Rollback()

		config.Logger.Fatal(err.Error())
	}

	time := time.Now().Format(time.UnixDate)
	// CheckParent 1 раз
	var checked bool
	for _, post := range posts {

		if !checked {
			if ok := CheckParent(post.Parent, treadID); !ok {
				status = http.StatusConflict
				tx.Rollback()
				postsId = []int{}

				return
			}

			checked = true
		}

		if ok := CheckUser(post.Author); !ok {
			status = http.StatusNotFound
			tx.Rollback()
			postsId = []int{}

			return
		}

		post.Thread = uint(treadID)
		post.Forum, _ = GetForumByThread(post.Thread)

		var postId int
		err := stmnt.QueryRow(
			post.Author,
			time,
			post.Forum,
			post.Message,
			post.Parent,
			post.Thread).Scan(&postId)
		if err != nil {
			tx.Rollback()
			postsId = []int{}

			config.Logger.Fatal("CreatePost", err.Error())
		}

		postsId = append(postsId, postId)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		postsId = []int{}

		config.Logger.Fatal("CreatePost", err.Error())
	}

	status = http.StatusCreated

	return
}

func SelectCreatedPosts(postsId []int) []api.Post {
	postsFull := []api.Post{}

	for _, postId := range postsId {
		postFull := api.Post{}

		row := config.DB.QueryRow(sqlSelectAnotherPostParams, postId)
		err := row.Scan(&postFull.Author,
			&postFull.Created,
			&postFull.Forum,
			&postFull.Id,
			&postFull.Message,
			&postFull.Parent,
			&postFull.Thread)
		if err != nil {
			log.Fatalln("SelectCreatedPosts, more params:", err.Error())
		}

		postsFull = append(postsFull, postFull)
	}

	return postsFull
}

func CheckParent(idParent api.JsonNullInt64, thread int) bool {
	if !idParent.Valid {
		return true
	}

	row := config.DB.QueryRow(sqlCheckParentPost, idParent, thread)
	err := row.Scan()
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func CheckPost(id int) bool {
	per := 0
	row := config.DB.QueryRow(sqlCheckPost, id)
	err := row.Scan(&per)
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

func UpdatePost(id int, updaet api.PostUpdaet) (status int) {
	if ok := CheckPost(id); !ok {
		return http.StatusNotFound
	}

	if updaet.Message == "" {
		return http.StatusOK
	}

	post, _ := SelectPost(id)

	if updaet.Message == post.Message {
		return http.StatusOK
	}

	res, err := config.DB.Exec(sqlPostUpdate, updaet.Message, id)
	if err != nil {
		log.Fatalln("UpdatePost", err.Error())
	}

	rows, _ := res.RowsAffected()
	if rows != 1 {
		log.Fatalf("UpdatePost: expected %v, have %v", 1, rows)
	}

	return http.StatusOK
}

func SelectPost(id int) (post api.Post, status int) {
	row := config.DB.QueryRow(sqlSelectPost, id)
	err := row.Scan(&post.Author,
		&post.Created,
		&post.Forum,
		&post.Id,
		&post.Message,
		&post.Parent,
		&post.Thread,
		&post.IsEdited)
	if err == sql.ErrNoRows {
		return post, http.StatusNotFound
	}

	return post, http.StatusOK
}

func SortedPosts(params *api.PostsSorted, thread int) []api.Post {
	var posts []api.Post
	if params.Sort == "flat" || params.Sort == "" {
		posts = flatSortedPosts(params, thread)
	} else if params.Sort == "tree" {
		posts = treeSortedPosts(params, thread)
	} else if params.Sort == "parent_tree" {
		posts = parentTreeSortedPosts(params, thread)
	}

	return posts
}

var sqlPosts = `select author, created, forum, id, message, parent, thread, isedited
from posts `

var sqlInsertPost = `INSERT INTO posts (author, created, forum, message, parent, thread, path)
 				VALUES ($1, $2, $3, $4, $5, $6,
                (SELECT path FROM posts WHERE id = $5)
                ||
				(SELECT currval('posts_id_seq')))
				RETURNING id`

var sqlCheckParentPost = `select id from posts where id = $1 and thread = $2`

var sqlSelectAnotherPostParams = `select author, created, forum, id, message, parent, thread
from posts where id = $1`

var sqlCheckPost = `select id from posts where id = $1`

var sqlPostUpdate = `update posts set message = (case when $1 = '' then message else $1 end), isedited = true
where id = $2`

var sqlSelectPost = `select author, created, forum, id, message, parent, thread, isedited
from posts where id = $1`
