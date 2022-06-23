package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/votes/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/sirupsen/logrus"
)

const (
	queryGetVoice      = `SELECT voice FROM votes WHERE nickname = $1 AND thread = $2;`
	queryCreate        = `INSERT INTO votes (nickname, thread, voice) VALUES ($1, $2, $3);`
	querySetByThreadId = `
							INSERT INTO votes (nickname, thread, voice) VALUES ($1, $2, $3)
					 		ON CONFLICT (nickname, thread) DO UPDATE
								SET voice = $3;`
	querySetByThreadSlug = `
							INSERT INTO votes (nickname, thread, voice) 
								SELECT $1, id, $3
								FROM threads
								WHERE slug = $2
					 		ON CONFLICT (nickname, thread) DO UPDATE
								SET voice = $3;`
	queryPatch = `
					UPDATE votes
					SET voice = COALESCE(NULLIF(TRIM($3), ''), voice)
					WHERE nickname = $1 AND thread = $2
					RETURNING voice;`
)

type VoteRepositoryPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *VoteRepositoryPostgres {
	return &VoteRepositoryPostgres{
		db: db,
	}
}

func (r *VoteRepositoryPostgres) Set(ctx context.Context, vote domain.Vote) (domain.Vote, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Vote",
		"method": "Set",
	})

	_, err := r.db.Exec(ctx, querySetByThreadId, vote.Nickname, vote.Thread, vote.Voice)
	if err != nil {
		log.Error(err.Error())

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.SQLState() {
			case "23503":
				return vote, forumErrors.NewEntityNotExistsError("votes")
			}
		}
	}

	return vote, err
}

func (r *VoteRepositoryPostgres) SetByThreadSlug(ctx context.Context, slug string, vote domain.Vote) (domain.Vote, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Vote",
		"method": "SetByThreadSlug",
	})

	_, err := r.db.Exec(ctx, querySetByThreadSlug, vote.Nickname, slug, vote.Voice)
	if err != nil {
		log.Error(err.Error())

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.SQLState() {
			case "23503":
				return vote, forumErrors.NewEntityNotExistsError("votes")
			}
		}
	}

	return vote, err
}

func (r *VoteRepositoryPostgres) Create(ctx context.Context, vote domain.Vote) (domain.Vote, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Vote",
		"method": "Create",
	})

	_, err := r.db.Exec(ctx, queryCreate, vote.Nickname, vote.Thread, vote.Voice)
	if err != nil {
		log.Error(err.Error())
	}

	return vote, err
}

func (r *VoteRepositoryPostgres) Exists(ctx context.Context, nickname string, thread int64) (bool, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Vote",
		"method": "Exists",
	})

	var voice int32

	err := r.db.QueryRow(ctx, queryGetVoice, nickname, thread).Scan(&voice)
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *VoteRepositoryPostgres) Patch(ctx context.Context, nickname string, thread int64, voice *int64) (domain.Vote, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Vote",
		"method": "Patch",
	})

	vote := domain.Vote{
		Nickname: nickname,
		Thread:   thread,
	}

	err := r.db.QueryRow(ctx, queryPatch, nickname, thread, voice).Scan(&vote.Voice)
	if err != nil {
		log.Error(err.Error())
	}

	return vote, err
}
