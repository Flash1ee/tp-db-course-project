package models

type ErrResponse struct {
	Err string `json:"error"`
}

type ResponseForum struct {
	Title         string `json:"title"`
	UsersNickname string `json:"user"`
	Slug          string `json:"slug"`
}
