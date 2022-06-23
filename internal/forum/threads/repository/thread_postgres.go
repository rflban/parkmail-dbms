package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	queryCreate = `INSERT INTO threads (title, author, forum, message, slug, created)
					VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id, title, author, forum, message, slug, created, votes;`
	queryCreate2 = `INSERT INTO threads (title, author, forum, message, slug)
					VALUES ($1, $2, $3, $4, $5)
					RETURNING id, title, author, forum, message, slug, created, votes;`
	queryGetById    = `SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1;`
	queryGetBySlug  = `SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE slug = $1;`
	queryUpdateById = `UPDATE threads SET
						title = COALESCE(NULLIF(TRIM($2), ''), title),
						message = COALESCE(NULLIF(TRIM($3), ''), message)
						WHERE id = $1
						RETURNING title, author, forum, message, votes, slug, created;`
	queryUpdateBySlug = `
							UPDATE threads
							SET
								 title = COALESCE(NULLIF(TRIM($2), ''), title),
								 message = COALESCE(NULLIF(TRIM($3), ''), message)
							WHERE slug = $1
							RETURNING id, title, author, forum, message, votes, created;`
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

	var (
		row  pgx.Row
		slug *string
	)

	if thread.Slug != "" {
		slug = &thread.Slug
	}

	var obtained domain.Thread

	if thread.Created.Equal(time.Time{}) {
		row = r.db.QueryRow(ctx, queryCreate2,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			slug,
		)
	} else {
		row = r.db.QueryRow(ctx, queryCreate,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			slug,
			thread.Created,
		)
	}

	var fetchedSlug *string = nil

	err := row.Scan(
		&obtained.Id,
		&obtained.Title,
		&obtained.Author,
		&obtained.Forum,
		&obtained.Message,
		&fetchedSlug,
		&obtained.Created,
		&obtained.Votes,
	)

	if fetchedSlug != nil {
		obtained.Slug = *fetchedSlug
	}

	if err != nil {
		log.Error(err.Error())

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.SQLState() {
			case "23505":
				return obtained, forumErrors.NewUniqueError(
					pgErr.TableName,
					pgErr.ColumnName,
				)
			case "23503":
				return obtained, forumErrors.NewEntityNotExistsError("users or forum")
			}
		}
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			err = forumErrors.NewUniqueError(
				pgErr.TableName,
				pgErr.ColumnName,
			)
		}
	}

	return obtained, err
}

func (r *ThreadRepositoryPostgres) GetById(ctx context.Context, id int64) (domain.Thread, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Thread",
		"method": "GetById",
	})

	var (
		thread domain.Thread
		slug   *string
	)

	err := r.db.QueryRow(ctx, queryGetById, id).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&slug,
		&thread.Created,
	)

	if slug != nil {
		thread.Slug = *slug
	}

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

	var (
		thread      domain.Thread
		fetchedSlug *string
	)

	err := r.db.QueryRow(ctx, queryGetBySlug, slug).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&fetchedSlug,
		&thread.Created,
	)

	if fetchedSlug != nil {
		thread.Slug = *fetchedSlug
	}

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

func (r *ThreadRepositoryPostgres) PatchBySlug(ctx context.Context, slug string, partialThread domain.PartialThread) (domain.Thread, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Thread",
		"method": "Patch",
	})

	thread := domain.Thread{
		Slug: slug,
	}
	err := r.db.QueryRow(ctx, queryUpdateById, slug, partialThread.Title, partialThread.Message).Scan(
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
