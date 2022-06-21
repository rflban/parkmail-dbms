package domain

import "github.com/rflban/parkmail-dbms/pkg/forum/models"

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

func (thread Thread) ToModel() models.Thread {
	id := int32(thread.Id)

	return models.Thread{
		Id:      &id,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   &thread.Forum,
		Message: thread.Message,
		Votes:   &thread.Votes,
		Slug:    &thread.Slug,
		Created: &thread.Created,
	}
}

func FromModel(thread models.Thread, id *int64) Thread {
	var (
		idVal      int64
		votesVal   int32
		forumVal   string
		slugVal    string
		createdVal string
	)

	if id != nil {
		idVal = *id
	}
	if thread.Votes != nil {
		votesVal = *thread.Votes
	}
	if thread.Forum != nil {
		forumVal = *thread.Forum
	}
	if thread.Slug != nil {
		slugVal = *thread.Slug
	}
	if thread.Created != nil {
		createdVal = *thread.Created
	}

	return Thread{
		Id:      idVal,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   forumVal,
		Message: thread.Message,
		Votes:   votesVal,
		Slug:    slugVal,
		Created: createdVal,
	}
}
