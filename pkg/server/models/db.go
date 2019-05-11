package models

import (
	"DB_Project_TP/api"
	"DB_Project_TP/config"
)

func GetDBInfo() api.DBInfo {
	info := api.DBInfo{}
	row := config.DB.QueryRow(sqlDBInfo)
	err := row.Scan(&info.Posts,
		&info.Threads,
		&info.Forums,
		&info.Users)
	if err != nil {
		config.Logger.Fatal("GetDBInfo", err.Error())
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
		config.Logger.Fatal("TruncateUserTable: Can not do pre exectute", err.Error())
	}
}

func TruncateThreadsTable() {
	_, err := config.DB.Exec(sqlTruncateThreads)
	if err != nil {
		config.Logger.Fatal("TruncateThreadsTable: Can not do pre exectute", err.Error())
	}
}

func TruncateForumsTable() {
	_, err := config.DB.Exec(sqlTruncateForums)
	if err != nil {
		config.Logger.Fatal("TruncateForumsTable: Can not do pre exectute", err.Error())
	}
}

func TruncatePostsTable() {
	_, err := config.DB.Exec(sqlTruncatePosts)
	if err != nil {
		config.Logger.Fatal("TruncatePostsTable: Can not do pre exectute", err.Error())
	}
}

var sqlTruncateUsers = `truncate table users CASCADE`

var sqlTruncateThreads = `truncate table threads CASCADE`

var sqlTruncateForums = `truncate table forums CASCADE`

var sqlTruncatePosts = `truncate table posts CASCADE`

var sqlDBInfo = `select
  (select count(*) from posts) as posts,
  (select count(*) from threads) as threads,
  (select count(*) from forums) as forums,
  (select count(*) from users) as users`
