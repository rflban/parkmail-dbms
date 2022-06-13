package domain

import "github.com/rflban/parkmail-dbms/pkg/forum/models"

type User struct {
	Nickname string
	Fullname string
	About    *string
	Email    string
}

func GetUserEntity(dto models.User) User {
	return User{
		Nickname: *dto.Nickname,
		Fullname: dto.Fullname,
		About:    dto.About,
		Email:    dto.Email,
	}
}

func (entity User) ToModel() models.User {
	return models.User{
		Nickname: &entity.Nickname,
		Fullname: entity.Fullname,
		About:    entity.About,
		Email:    entity.Email,
	}
}
