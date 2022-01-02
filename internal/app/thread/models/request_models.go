package models

type RequestUpdateThread struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RequestNewPost struct {
	Parent  int64  `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}
