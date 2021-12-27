package models

type RequestCreateForum struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type RequestCreateThread struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Created string `json:"created"`
}
