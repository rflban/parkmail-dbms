package models

//easyjson:json
type UserUpdate struct {
	Fullname *string `json:"fullname,omitempty"`
	About    *string `json:"about,omitempty"`
	Email    *string `json:"email,omitempty"`
}
