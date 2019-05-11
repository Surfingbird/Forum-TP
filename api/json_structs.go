package api

import (
	"database/sql"
	"encoding/json"
)

type Forum struct {
	Posts   uint64 `json:"posts"`
	Slug    string `json:"slug"`
	Threads uint64 `json:"threads"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

type ForumsUsersQuery struct {
	Limit int    `schema:"limit"`
	Since string `schema:"since"`
	Desc  bool   `schema:"desc"`
}

type DBInfo struct {
	Forums  uint `json:"forum"`
	Posts   uint `json:"post"`
	Threads uint `json:"thread"`
	Users   uint `json:"user"`
}

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type Post struct {
	Author   string        `json:"author"`
	Created  string        `json:"created"`
	Forum    string        `json:"forum"`
	Id       uint          `json:"id"`
	IsEdited bool          `json:"isEdited"`
	Message  string        `json:"message"`
	Parent   JsonNullInt64 `json:"parent"`
	Thread   uint          `json:"thread"`
}

type PostUpdaet struct {
	Message string `json:"message"`
}

type PostsSorted struct {
	Limit int    `schema:"limit"`
	Since string `schema:"since"`
	Desc  bool   `schema:"desc"`
	Sort  string `schems:"sort"`
}

type PostParams struct {
	Related []string `schema:"related"`
}

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
}

type UpdateUser struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      uint64 `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Votes   int64  `json:"votes"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type FullPost struct {
	User
	Forum
	Post
	Thread
}

type Error struct {
	Message string `json:"message"`
}
