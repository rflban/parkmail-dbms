package errors

import "fmt"

type ConflictError struct {
	message string
}

func NewConflictError(message string) ConflictError {
	return ConflictError{
		message: message,
	}
}

func (e ConflictError) Error() string {
	return e.message
}

type UniqueError struct {
	entity string
	attr   string
}

func NewUniqueError(entity, attr string) UniqueError {
	return UniqueError{
		entity: entity,
		attr:   attr,
	}
}

func (e UniqueError) Error() string {
	return fmt.Sprintf("Unique constraint violation at '%s' entity '%s' attribute", e.entity, e.attr)
}

type EntityNotExistsError struct {
	entity string
}

func NewEntityNotExistsError(entity string) EntityNotExistsError {
	return EntityNotExistsError{
		entity: entity,
	}
}

func (e EntityNotExistsError) Error() string {
	return fmt.Sprintf("Not found for entity '%s'", e.entity)
}
