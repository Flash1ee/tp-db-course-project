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
	Forum   string `json:"forum,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Created string `json:"created"`
}
