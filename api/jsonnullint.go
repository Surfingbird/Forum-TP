package api

import (
	"time"
)

type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	Id       uint      `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parent"`
	Thread   uint      `json:"thread"`
	Root     int64     `json:"-"`
}
