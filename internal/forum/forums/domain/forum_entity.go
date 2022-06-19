package domain

type Forum struct {
	Id      int64
	Title   string
	User    string
	Slug    string
	Posts   int64
	Threads int32
}
