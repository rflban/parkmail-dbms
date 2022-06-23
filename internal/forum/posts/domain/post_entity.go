package domain

import (
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"time"
)

type Post struct {
	Id       int64
	Parent   int64
	Author   string
	Message  string
	IsEdited bool
	Forum    string
	Thread   int64
	Created  time.Time
}

func (post Post) ToModel() models.Post {
	thread := int32(post.Thread)

	return models.Post{
		Id:       &post.Id,
		Parent:   &post.Parent,
		Author:   post.Author,
		Message:  post.Message,
		IsEdited: &post.IsEdited,
		Forum:    &post.Forum,
		Thread:   &thread,
		Created:  &post.Created,
	}
}

func FromModel(post models.Post) Post {
	var (
		idVal       int64
		parentVal   int64
		isEditedVal bool
		forumVal    string
		threadVal   int64
		createdVal  time.Time
	)

	if post.Id != nil {
		idVal = *post.Id
	}
	if post.Parent != nil {
		parentVal = *post.Parent
	}
	if post.IsEdited != nil {
		isEditedVal = *post.IsEdited
	}
	if post.Forum != nil {
		forumVal = *post.Forum
	}
	if post.Thread != nil {
		threadVal = int64(*post.Thread)
	}
	if post.Created != nil {
		createdVal = *post.Created
	}

	return Post{
		Id:       idVal,
		Parent:   parentVal,
		Author:   post.Author,
		Message:  post.Message,
		IsEdited: isEditedVal,
		Forum:    forumVal,
		Thread:   threadVal,
		Created:  createdVal,
	}
}
