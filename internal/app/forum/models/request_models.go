package models

type RequestCreateForum struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}
