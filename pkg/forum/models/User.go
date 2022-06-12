package models

//easyjson:json
type User struct {
	Nickname *string `json:"nickname,omitempty"`
	Fullname string  `json:"fullname"`
	About    *string `json:"about,omitempty"`
	Email    string  `json:"email"`
}
