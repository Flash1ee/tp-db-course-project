package models

type Forum struct {
	Title         string `json:"title"`
	UsersNickname string `json:"user"`
	Slug          string `json:"slug"`
	Posts         int64  `json:"posts"`
	Threads       int64  `json:"threads"`
}
