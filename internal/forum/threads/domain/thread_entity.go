package domain

type Thread struct {
	Id      int64
	Title   string
	Author  string
	Forum   string
	Message string
	Votes   int32
	Slug    string
	Created string
}
