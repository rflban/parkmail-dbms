package models

//easyjson:json
type ThreadUpdate struct {
	Title   *string `json:"title,omitempty"`
	Message *string `json:"message,omitempty"`
}
