package domain

type Post struct {
	Id       int64
	Parent   int64
	Author   string
	Message  string
	IsEdited bool
	Forum    string
	Thread   int64
	Created  string
}
