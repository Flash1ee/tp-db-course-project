package models

type RequestVoteUpdate struct {
	Nickname string `json:"nickname"`
	Voice    int64  `json:"voice"`
}
