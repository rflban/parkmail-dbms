package domain

import "github.com/rflban/parkmail-dbms/pkg/forum/models"

type PartialThread struct {
	Title   *string
	Message *string
}

func FromModelUpdate(thread models.ThreadUpdate) PartialThread {
	return PartialThread{
		Title:   thread.Title,
		Message: thread.Message,
	}
}
