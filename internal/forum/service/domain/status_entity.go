package domain

import (
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type Status struct {
	User   int32
	Forum  int32
	Thread int32
	Post   int64
}

func (entity Status) ToModel() models.Status {
	return models.Status{
		User:   entity.User,
		Forum:  entity.Forum,
		Thread: entity.Thread,
		Post:   entity.Post,
	}
}
