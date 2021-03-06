package models

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"DB_Project_TP/api"
	"DB_Project_TP/config"
)

func CreateUser(u *api.User) (status int) {
	_, err := config.DB.Exec(sqlInsertUser,
		u.Fullname,
		u.Nickname,
		u.Email,
		u.About)
	if err != nil {
		log.Println("CreateUser", err.Error(), "expected StatusConflict")

		return http.StatusConflict
	}

	return http.StatusCreated
}

func SelectConflictUsers(nickname, email string) []api.User {
	rows, err := config.DB.Query(sqlSelectUserConflict, nickname, email)
	if err != nil {
		log.Fatalln("SelectConflictUsers", err.Error())
	}
	defer rows.Close()

	users := []api.User{}

	for rows.Next() {
		user := api.User{}

		if err := rows.Scan(&user.About,
			&user.Email,
			&user.Fullname,
			&user.Nickname); err != nil {
			log.Fatalf("ProfileHandler: %v\n", err.Error())
		}

		users = append(users, user)
	}

	return users
}

func SelectUser(nickname string) (u api.User, err error) {
	err = config.DB.QueryRow(sqlSelectUser, nickname).Scan(&u.About,
		&u.Email,
		&u.Fullname,
		&u.Nickname)
	if err != nil {
		err = errors.New("There is no this user!")

		return u, err
	}

	return u, nil
}

func CheckUser(nickname string) bool {
	result, err := config.DB.Exec(sqlCheckUser, nickname)
	if err != nil {
		log.Fatalln("CheckUser", err.Error())
	}

	rows := result.RowsAffected()
	if rows != 1 {
		return false
	}

	return true
}

func UpdateUser(update *api.UpdateUser, nickname string) (u api.User, status int) {
	if ok := CheckUser(nickname); !ok {
		return u, http.StatusNotFound
	}

	result, err := config.DB.Exec(sqlUpdateUser,
		update.About,
		update.Email,
		update.Fullname,
		nickname)
	if err != nil {
		return u, http.StatusConflict
	}

	row := result.RowsAffected()
	if row != 1 {
		log.Fatalln("Count of updated != 1")
	}

	u, err = SelectUser(nickname)
	if err != nil {
		log.Fatalln("Can not select updated user!")
	}

	return u, http.StatusOK
}

func RegNickname(nickname string) (regNick string, status int) {
	row, err := config.DB.Query(sqlRegNickname, nickname)
	if err != nil {
		log.Fatalln("RegNickname", err.Error())
	}
	defer row.Close()

	if !row.Next() {
		return "", http.StatusNotFound
	}

	if err := row.Scan(&regNick); err != nil {
		log.Fatalf("RegNickname: %v\n", err.Error())
	}

	return regNick, http.StatusOK
}

func QueryForumUsersWithParams(params api.ForumsUsersQuery) string {
	sort := " ASC "
	limit := " ALL "
	compare := " > "
	query := sqlForumsUsers
	if params.Desc {
		sort = " DESC "
		compare = " < "
	}
	if params.Limit != 0 {
		limit = strconv.Itoa(params.Limit)
	}
	if params.Since != "" {
		query += ` and u.nickname ` + compare + "'" + params.Since + "'"
	}

	query += ` order by lower(u.nickname COLLATE "C") ` + sort + " limit " + limit

	return query
}

func SelectForumsUsers(params api.ForumsUsersQuery, slug string) (users []api.User, status int) {
	users = []api.User{}

	if ok, _ := CheckForum(slug); !ok {
		return users, http.StatusNotFound
	}

	query := QueryForumUsersWithParams(params)

	rows, err := config.DB.Query(query, slug)
	if err != nil {
		log.Fatalln("SelectPosts", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		user := api.User{}
		err := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			log.Fatalln("SelectForumsUsers", err.Error())
		}

		users = append(users, user)
	}

	return users, http.StatusOK
}

var sqlInsertUser = `INSERT INTO users 
							(fullname, nickname, email, about)
							VALUES ($1, $2, $3, $4)`

var sqlSelectUser = `select
					about,
					email,
					fullname,
					nickname
				from users
				where nickname = $1`

var sqlSelectUserConflict = `select
  							about,
							email,
  							fullname,
  							nickname
						from users
						where nickname = $1 
						or email = $2`

var sqlUpdateUser = `update users
set about = (case
            when $1 = '' then about
             else $1 end),
    email = (case
            when $2 = '' then email
            else $2 end),
    fullname = (case
            when $3 = '' then fullname 
            else $3 end)
where nickname = $4`

var sqlCheckUser = `select * from users where nickname = $1`

var sqlRegNickname = `select nickname  from users where nickname = $1`

var sqlForumsUsers = `
select u.about, u.email, u.fullname, u.nickname
from users u

where (
  exists(
      select *
      from posts p
      where p.author = u.nickname and p.forum = $1
  )
  or
  exists(
      select *
      from threads t
      where t.author = u.nickname and t.forum = $1
  )
) 
`
