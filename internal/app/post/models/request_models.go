package models

type RequestUpdateMessage struct {
	Message string `json:"message"`
}
type RequestCreatePost struct {
	Parent  uint64 `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}
