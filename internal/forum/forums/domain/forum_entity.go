package domain

import "github.com/rflban/parkmail-dbms/pkg/forum/models"

type Forum struct {
	Id      int64
	Title   string
	User    string
	Slug    string
	Posts   int64
	Threads int32
}

func (forum Forum) ToModel() models.Forum {
	return models.Forum{
		Title:   forum.Title,
		User:    forum.User,
		Slug:    forum.Slug,
		Posts:   &forum.Posts,
		Threads: &forum.Threads,
	}
}

func FromModel(forum models.Forum, id *int64) Forum {
	var (
		idVal      int64
		postsVal   int64
		threadsVal int32
	)

	if id != nil {
		idVal = *id
	}
	if forum.Posts != nil {
		postsVal = *forum.Posts
	}
	if forum.Threads != nil {
		threadsVal = *forum.Threads
	}

	return Forum{
		Id:      idVal,
		Title:   forum.Title,
		User:    forum.User,
		Slug:    forum.Slug,
		Posts:   postsVal,
		Threads: threadsVal,
	}
}
