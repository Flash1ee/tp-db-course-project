package models

import (
	"time"
	models_forum "tp-db-project/internal/app/forum/models"
	models_thread "tp-db-project/internal/app/thread/models"
	models_author "tp-db-project/internal/app/users/models"
)

type ResponsePost struct {
	Id       int64     `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int64     `json:"thread"`
	Created  time.Time `json:"created"`
}
type ResponsePostDetail struct {
	Post   *ResponsePost                 `json:"post"`
	Author *models_author.ResponseUser   `json:"author,omitempty"`
	Thread *models_thread.ResponseThread `json:"thread,omitempty"`
	Forum  *models_forum.ResponseForum   `json:"forum,omitempty"`
}
