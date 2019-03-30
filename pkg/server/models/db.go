package models

import (
	"DB_Project_TP/api"
	"DB_Project_TP/config"
	"log"
)

func GetDBInfo() api.DBInfo {
	info := api.DBInfo{}
	row := config.DB.QueryRow(sqlDBInfo)
	err := row.Scan(&info.Posts,
		&info.Threads,
		&info.Forums,
		&info.Users)
	if err != nil {
		log.Fatalln("GetDBInfo", err.Error())
	}

	return info
}

func TruncateAllTables() {
	TruncateUserTable()
	TruncateThreadsTable()
	TruncateForumsTable()
	TruncatePostsTable()
}

func TruncateUserTable() {
	_, err := config.DB.Exec(sqlTruncateUsers)
	if err != nil {
		log.Fatalln("Can not do pre exectute", err.Error())
	}
}

func TruncateThreadsTable() {
	_, err := config.DB.Exec(sqlTruncateThreads)
	if err != nil {
		log.Fatalln("Can not do pre exectute", err.Error())
	}
}

func TruncateForumsTable() {
	_, err := config.DB.Exec(sqlTruncateForums)
	if err != nil {
		log.Fatalln("Can not do pre exectute", err.Error())
	}
}

func TruncatePostsTable() {
	_, err := config.DB.Exec(sqlTruncatePosts)
	if err != nil {
		log.Fatalln("Can not do pre exectute", err.Error())
	}
}

var sqlTruncateUsers = `truncate table project_bd.users`

var sqlTruncateThreads = `truncate table project_bd.threads`

var sqlTruncateForums = `truncate table project_bd.forums`

var sqlTruncatePosts = `truncate table project_bd.posts`

var sqlDBInfo = `select
  (select count(*) from project_bd.posts) as posts,
  (select count(*) from project_bd.threads) as threads,
  (select count(*) from project_bd.forums) as forums,
  (select count(*) from project_bd.users) as users`
