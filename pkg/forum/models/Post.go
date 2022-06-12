package models

//easyjson:json
type Post struct {
	Id       *int64  `json:"id,omitempty"`
	Parent   *int64  `json:"parent,omitempty"`
	Author   string  `json:"author"`
	Message  string  `json:"message"`
	IsEdited *bool   `json:"isEdited,omitempty"`
	Forum    *string `json:"forum,omitempty"`
	Thread   *int32  `json:"thread,omitempty"`
	Created  *string `json:"created,omitempty"`
}
