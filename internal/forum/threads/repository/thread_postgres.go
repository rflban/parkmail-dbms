package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/sirupsen/logrus"
)

const (
	queryCreate = `INSERT INTO threads (title, author, forum, message, slug, created)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id, votes;`
	queryGetById    = `SELECT title, author, forum, message, votes, slug, created FROM threads WHERE id = $1;`
	queryGetBySlug  = `SELECT id, title, author, forum, message, votes, created FROM threads WHERE id = $1;`
	queryUpdateById = `UPDATE threads SET
						title = COALESCE(NULLIF(TRIM($2), ''), title),
						message = COALESCE(NULLIF(TRIM($3), ''), message)
						WHERE id = $1
						RETURNING title, author, forum, message, votes, slug, created;`
)

type ThreadRepositoryPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *ThreadRepositoryPostgres {
	return &ThreadRepositoryPostgres{
		db: db,
	}
}

func (r *ThreadRepositoryPostgres) Create(ctx context.Context, thread domain.Thread) (domain.Thread, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Thread",
		"method": "Create",
	})

	err := r.db.QueryRow(ctx, queryCreate,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
		thread.Slug,
		thread.Created,
	).Scan(&thread.Id, &thread.Votes)

	if err != nil {
		log.Error(err.Error())

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			err = forumErrors.NewUniqueError(
				pgErr.TableName,
				pgErr.ColumnName,
			)
		}
	}

	return thread, err
}

func (r *ThreadRepositoryPostgres) GetById(ctx context.Context, id int64) (domain.Thread, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Thread",
		"method": "GetById",
	})

	thread := domain.Thread{
		Id: id,
	}
	err := r.db.QueryRow(ctx, queryGetById, id).Scan(
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		log.Error(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return thread, forumErrors.NewEntityNotExistsError("threads")
		}
	}

	return thread, err
}

func (r *ThreadRepositoryPostgres) GetBySlug(ctx context.Context, slug string) (domain.Thread, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Thread",
		"method": "GetBySlug",
	})

	thread := domain.Thread{
		Slug: slug,
	}
	err := r.db.QueryRow(ctx, queryGetBySlug, slug).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Created,
	)

	if err != nil {
		log.Error(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return thread, forumErrors.NewEntityNotExistsError("threads")
		}
	}

	return thread, err
}

func (r *ThreadRepositoryPostgres) Patch(ctx context.Context, id int64, partialThread domain.PartialThread) (domain.Thread, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Thread",
		"method": "Patch",
	})

	thread := domain.Thread{
		Id: id,
	}
	err := r.db.QueryRow(ctx, queryUpdateById, id, partialThread.Title, partialThread.Message).Scan(
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		log.Error(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return thread, forumErrors.NewEntityNotExistsError("threads")
		}
	}

	return thread, err
}
