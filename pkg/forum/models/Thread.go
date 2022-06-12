package models

type Thread struct {
	Id      *int32  `json:"id,omitempty"`
	Title   string  `json:"title"`
	Author  string  `json:"author"`
	Forum   *string `json:"forum,omitempty"`
	Message string  `json:"message"`
	Votes   *int32  `json:"votes,omitempty"`
	Slug    *string `json:"slug,omitempty"`
	Created *string `json:"created,omitempty"`
}
