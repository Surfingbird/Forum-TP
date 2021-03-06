package models

import (
	"log"
	"net/http"
	"time"

	"DB_Project_TP/api"
	"DB_Project_TP/config"
)

func CreatePost(posts []*api.Post, treadID int) (status int) {

	tx, _ := config.DB.Begin()
	_, err := tx.Prepare("bulk_create", sqlInsertPost)
	if err != nil {
		tx.Rollback()

		config.Logger.Fatal(err.Error())
	}

	time := time.Now()

	if len(posts) > 0 {
		if ok := CheckParent(posts[0].Parent, treadID); !ok {
			status = http.StatusConflict
			tx.Rollback()

			return
		}

		if ok := CheckUser(posts[0].Author); !ok {
			status = http.StatusNotFound
			tx.Rollback()

			return
		}
	}

	for _, post := range posts {
		post.Thread = uint(treadID)
		post.Forum, _ = GetForumByThread(post.Thread)

		err := tx.QueryRow(
			"bulk_create",
			time,
			post.Author,
			post.Forum,
			post.Message,
			post.Parent,
			post.Thread).Scan(&post.Id, &post.Created)
		if err != nil {
			tx.Rollback()

			config.Logger.Fatal("CreatePost ", err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		config.Logger.Fatal("CreatePost ", err.Error())
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

func CheckParent(idParent int64, thread int) bool {
	if idParent == 0 {
		return true
	}

	var id int
	err := config.DB.QueryRow(sqlCheckParentPost, idParent, thread).Scan(&id)
	if err != nil {
		config.Logger.Info("CheckParent", err.Error())
		return false
	}

	return true
}

func CheckPost(id int) bool {
	per := 0
	err := config.DB.QueryRow(sqlCheckPost, id).Scan(&per)
	if err != nil {
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

	rows := res.RowsAffected()
	if rows != 1 {
		log.Fatalf("UpdatePost: expected %v, have %v", 1, rows)
	}

	return http.StatusOK
}

func SelectPost(id int) (post api.Post, status int) {
	err := config.DB.QueryRow(sqlSelectPost, id).Scan(&post.Author,
		&post.Created,
		&post.Forum,
		&post.Id,
		&post.Message,
		&post.Parent,
		&post.Thread,
		&post.IsEdited)
	if err != nil {
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

var sqlInsertPost = ` INSERT INTO posts (created, author, forum, message, parent, thread, path, post_root) VALUES ($1, $2, $3, $4, $5, $6,
	(SELECT path FROM posts WHERE id = $5)
			||
			(SELECT currval('posts_id_seq')),
			CASE WHEN $5 = 0
				THEN currval('posts_id_seq')
				ELSE
					(SELECT post_root FROM posts WHERE id = $5)
			END)

			RETURNING id, created`

var sqlCheckParentPost = `select id from posts where id = $1 and thread = $2`

var sqlSelectAnotherPostParams = `select author, created, forum, id, message, parent, thread
from posts where id = $1`

var sqlCheckPost = `select id from posts where id = $1`

var sqlPostUpdate = `update posts set message = (case when $1 = '' then message else $1 end), isedited = true
where id = $2`

var sqlSelectPost = `select author, created, forum, id, message, parent, thread, isedited
from posts where id = $1`
