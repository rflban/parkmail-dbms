package domain

import "github.com/rflban/parkmail-dbms/pkg/forum/models"

type Vote struct {
	Thread   int64
	Nickname string
	Voice    int32
}

func (vote Vote) ToModel() models.Vote {
	return models.Vote{
		Nickname: vote.Nickname,
		Voice:    vote.Voice,
	}
}

func FromModel(vote models.Vote, thread int64) Vote {
	return Vote{
		Thread:   thread,
		Nickname: vote.Nickname,
		Voice:    vote.Voice,
	}
}
