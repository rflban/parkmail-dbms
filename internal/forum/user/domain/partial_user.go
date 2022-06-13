package domain

import "github.com/rflban/parkmail-dbms/pkg/forum/models"

type PartialUser struct {
	Fullname *string
	About    *string
	Email    *string
}

func GetPartial(dto models.UserUpdate) PartialUser {
	return PartialUser{
		Fullname: dto.Fullname,
		About:    dto.About,
		Email:    dto.Email,
	}
}
