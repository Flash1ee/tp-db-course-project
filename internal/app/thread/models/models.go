package models

import "time"

type Thread struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int64  `json:"votes"`
	Slug    string `json:"slug,omitempty"`
	Created time.Time `json:"created,omitempty"`
}
