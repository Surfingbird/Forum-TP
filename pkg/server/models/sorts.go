package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx"

	"DB_Project_TP/api"

	"DB_Project_TP/config"
)

func flatSortedPosts(params *api.PostsSorted, thread int) []api.Post {
	sort := " ASC "
	limit := " ALL "
	compare := " > "

	if params.Limit != 0 {
		limit = strconv.Itoa(params.Limit)
	}
	if params.Desc {
		sort = " DESC "
		compare = " < "
	}

	where := fmt.Sprintf(" where thread = $1 ")
	orderBY := fmt.Sprintf(" order by created %v, id %v", sort, sort)
	limitStr := fmt.Sprintf(" limit %v ", limit)

	var rows *pgx.Rows
	var err error
	if params.Since != "" {
		where += fmt.Sprintf(" and id %v $2 ", compare)
		query := sqlPosts + where + orderBY + limitStr

		rows, err = config.DB.Query(query, thread, params.Since)
		if err != nil {
			log.Fatalln("flatSortedPosts", err.Error())
		}
	} else {
		query := sqlPosts + where + orderBY + limitStr

		rows, err = config.DB.Query(query, thread)
		if err != nil {
			log.Fatalln("flatSortedPosts", err.Error())
		}
	}

	posts := []api.Post{}
	for rows.Next() {
		post := api.Post{}
		err := rows.Scan(&post.Author,
			&post.Created,
			&post.Forum,
			&post.Id,
			&post.Message,
			&post.Parent,
			&post.Thread,
			&post.IsEdited)
		if err != nil {
			log.Fatalln("flatSortedPosts", err.Error())
		}

		posts = append(posts, post)
	}

	return posts
}

func treeSortedPosts(params *api.PostsSorted, thread int) []api.Post {
	sort := " ASC "
	limit := " ALL "
	compare := " > "

	if params.Limit != 0 {
		limit = strconv.Itoa(params.Limit)
	}
	if params.Desc {
		sort = " DESC "
		compare = " < "
	}

	where := fmt.Sprintf(" where thread = $1 ")
	orderBY := fmt.Sprintf(" order by path %v ", sort)
	limitStr := fmt.Sprintf(" limit %v ", limit)

	var rows *pgx.Rows
	var err error
	if params.Since != "" {
		where += fmt.Sprintf(" and path %v (SELECT path FROM posts p WHERE p.id = $2) ", compare)
		query := sqlPosts + where + orderBY + limitStr

		rows, err = config.DB.Query(query, thread, params.Since)
		if err != nil {
			log.Fatalln("treeSortedPosts", err.Error())
		}
	} else {
		query := sqlPosts + where + orderBY + limitStr

		rows, err = config.DB.Query(query, thread)
		if err != nil {
			log.Fatalln("treeSortedPosts", err.Error())
		}
	}

	posts := []api.Post{}
	for rows.Next() {
		post := api.Post{}
		err := rows.Scan(&post.Author,
			&post.Created,
			&post.Forum,
			&post.Id,
			&post.Message,
			&post.Parent,
			&post.Thread,
			&post.IsEdited)
		if err != nil {
			log.Fatalln("flatSortedPosts", err.Error())
		}

		posts = append(posts, post)
	}

	return posts
}

func parentTreeSortedPosts(params *api.PostsSorted, thread int) []api.Post {
	// sort := " ASC "
	// limit := " ALL "
	// compare := " > "

	// if params.Limit != 0 {
	// 	limit = strconv.Itoa(params.Limit)
	// }

	// orderBY := fmt.Sprintf(" order by path")
	// if params.Desc {
	// 	sort = " DESC "
	// 	orderBY = fmt.Sprintf(" order by path[1] %v, path", sort)
	// 	compare = " < "
	// }

	// where := fmt.Sprintf(" where thread = $1 and path[1] in ")

	// selectIn := fmt.Sprintf(" select par.id from posts par ")
	// whereIn := fmt.Sprintf(" where par.thread = $1 and par.parent is null ")
	// if params.Since != "" {
	// 	whereIn += fmt.Sprintf(" and par.path[1] %v (select path[1] from posts where id = $2) ", compare)
	// }

	// orderIn := fmt.Sprintf(" order by par.created %v, par.id %v ", sort, sort)
	// limitIn := fmt.Sprintf(" limit %v ", limit)
	// subQuery := " ( " + selectIn + whereIn + orderIn + limitIn + " ) "

	// query := sqlPosts + where + subQuery + orderBY

	// var rows *pgx.Rows
	// var err error
	// if params.Since != "" {
	// 	rows, err = config.DB.Query(query, thread, params.Since)
	// 	if err != nil {
	// 		log.Fatalln("parentTreeSortedPosts", err.Error())
	// 	}
	// } else {
	// 	rows, err = config.DB.Query(query, thread)
	// 	if err != nil {
	// 		log.Fatalln("parentTreeSortedPosts", err.Error())
	// 	}
	// }

	// posts := []api.Post{}
	// for rows.Next() {
	// 	post := api.Post{}
	// 	err := rows.Scan(&post.Author,
	// 		&post.Created,
	// 		&post.Forum,
	// 		&post.Id,
	// 		&post.Message,
	// 		&post.Parent,
	// 		&post.Thread,
	// 		&post.IsEdited)
	// 	if err != nil {
	// 		log.Fatalln("parentTreeSortedPosts", err.Error())
	// 	}

	// 	posts = append(posts, post)
	// }
	query := `
		SELECT id, parent, author, message, isedited, forum, thread, created 
		FROM posts 
	`

	query += ` WHERE post_root IN ( SELECT id FROM posts WHERE thread = $1 AND parent = 0 `
	eqOp := ""
	if params.Desc {
		eqOp = " < "
	} else {
		eqOp = " > "
	}

	if params.Since != "" {
		query += fmt.Sprintf(` AND id %s (SELECT post_root FROM posts WHERE id = %s) `, eqOp, params.Since)
	}

	sortOrd := ""
	sortOrd = ` ASC `
	if params.Desc {
		sortOrd = ` DESC `
	}

	query += fmt.Sprintf(` ORDER BY id %s `, sortOrd)

	if params.Limit != 0 {
		query += fmt.Sprintf(` LIMIT %v `, params.Limit)
	}

	query += `)`
	if params.Desc {
		query += ` ORDER BY post_root DESC, path `
	} else {
		query += ` ORDER BY path`
	}

	rows, err := config.DB.Query(query, thread)
	if err != nil {
		config.Logger.Fatal("parentTreeSortedPosts", err.Error())
	}
	defer rows.Close()

	posts := []api.Post{}

	for rows.Next() {
		post := api.Post{}

		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
			&post.IsEdited, &post.Forum, &post.Thread, &post.Created)

		posts = append(posts, post)
	}

	return posts
}
