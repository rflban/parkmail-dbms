package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/forums/domain"
	threadsDomain "github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	usersDomain "github.com/rflban/parkmail-dbms/internal/forum/users/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/sirupsen/logrus"
)

const (
	queryCreate    = `INSERT INTO forums (title, "user", slug, posts, threads) VALUES ($1, $2, $3, $4, $5);`
	queryGetBySlug = `SELECT id, title, "user", posts, threads FROM forums WHERE slug = $1;`
)

type ForumRepositoryPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *ForumRepositoryPostgres {
	return &ForumRepositoryPostgres{
		db: db,
	}
}

func (r *ForumRepositoryPostgres) Create(ctx context.Context, forum domain.Forum) (domain.Forum, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Forum",
		"method": "Create",
	})

	_, err := r.db.Exec(ctx, queryCreate, forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads)

	if err != nil {
		log.Error(err.Error())

		// TODO: user not found

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return forum, forumErrors.NewUniqueError(
				pgErr.TableName,
				pgErr.ColumnName,
			)
		}
	}

	return forum, err
}

func (r *ForumRepositoryPostgres) GetBySlug(ctx context.Context, slug string) (domain.Forum, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Forum",
		"method": "GetBySlug",
	})

	forum := domain.Forum{
		Slug: slug,
	}
	err := r.db.QueryRow(ctx, queryGetBySlug, slug).Scan(
		&forum.Id,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
	)

	if err != nil {
		log.Error(err.Error())

		if err.Error() == pgx.ErrNoRows.Error() {
			return forum, forumErrors.NewEntityNotExistsError("forums")
		}
	}

	return forum, err
}

func (r *ForumRepositoryPostgres) GetUsersBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) ([]usersDomain.User, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Forum",
		"method": "GetUsersBySlug",
	})

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("nickname, fullname, about, email").
		From("forums_users").
		Where("forum = ?", slug)

	if since != "" {
		if desc {
			queryBuilder = queryBuilder.Where(`nickname < ?`, since)
		} else {
			queryBuilder = queryBuilder.Where(`nickname > ?`, since)
		}
	}

	if desc {
		queryBuilder = queryBuilder.OrderBy(`nickname DESC`)
	} else {
		queryBuilder = queryBuilder.OrderBy(`nickname ASC`)
	}

	if limit > 0 {
		queryBuilder = queryBuilder.Limit(limit)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	users := make([]usersDomain.User, 0, rows.CommandTag().RowsAffected())
	user := usersDomain.User{}

	for rows.Next() {
		err = rows.Scan(
			&user.Id,
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *ForumRepositoryPostgres) GetThreadsBySlug(ctx context.Context, slug string, since string, limit uint64, desc bool) ([]threadsDomain.Thread, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Forum",
		"method": "GetThreadsBySlug",
	})

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id, title, author, message, votes, slug, created").
		From("threads").
		Where("forum = ?", slug)

	if since != "" {
		if desc {
			queryBuilder = queryBuilder.Where("created <= ?", since)
		} else {
			queryBuilder = queryBuilder.Where("created >= ?", since)
		}
	}

	if desc {
		queryBuilder.OrderBy("created DESC")
	} else {
		queryBuilder.OrderBy("created ASC")
	}

	if limit > 0 {
		queryBuilder.Limit(limit)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	threads := make([]threadsDomain.Thread, 0, rows.CommandTag().RowsAffected())
	thread := threadsDomain.Thread{
		Forum: slug,
	}

	for rows.Next() {
		err = rows.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}