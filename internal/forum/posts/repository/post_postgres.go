package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/posts/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/sirupsen/logrus"
)

const (
	queryLastId  = `SELECT MAX(id) FROM posts;`
	queryGetById = `SELECT parent, author, message, is_edited, forum, thread, created FROM posts WHERE id = $1;`
	queryUpdate  = `UPDATE posts
					SET message = COALESCE(NULLIF(TRIM($2), ''), message), is_edited = true
					WHERE id = $1
					RETURNING parent, author, message, is_edited, forum, thread, created;`
)

type PostRepositoryPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostRepositoryPostgres {
	return &PostRepositoryPostgres{
		db: db,
	}
}

func (r *PostRepositoryPostgres) Create(ctx context.Context, posts []domain.Post) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "Create",
	})

	var lastId int64
	err := r.db.QueryRow(ctx, queryLastId).Scan(&lastId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	copied, err := r.db.CopyFrom(ctx, pgx.Identifier{"posts"}, []string{
		"parent",
		"author",
		"message",
		"forum",
		"thread",
		"created",
	}, pgx.CopyFromSlice(len(posts), func(i int) ([]interface{}, error) {
		return []interface{}{
			posts[i].Parent,
			posts[i].Author,
			posts[i].Message,
			posts[i].Forum,
			posts[i].Thread,
			posts[i].Created,
		}, nil
	}))

	if err != nil {
		log.Error(err.Error())

		if err := tx.Rollback(ctx); err != nil {
			log.Error(err.Error())
		}

		return nil, err
	}

	if int(copied) != len(posts) {
		log.
			WithField("copied", fmt.Sprintf("%d/%d", copied, len(posts))).
			Errorf("Failed bulck insert")

		if err := tx.Rollback(ctx); err != nil {
			log.Error(err.Error())
		}

		return nil, err
	}

	for i := range posts {
		posts[i].Id = int64(i) + lastId
	}

	return posts, err
}

func (r *PostRepositoryPostgres) Patch(ctx context.Context, id int64, message *string) (domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "Patch",
	})

	post := domain.Post{
		Id: id,
	}
	err := r.db.QueryRow(ctx, queryUpdate, id, message).Scan(
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)

	if err != nil {
		log.Error(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return post, forumErrors.NewEntityNotExistsError("posts")
		}
	}

	return post, err
}

func (r *PostRepositoryPostgres) GetById(ctx context.Context, id int64) (domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetById",
	})

	post := domain.Post{
		Id: id,
	}
	err := r.db.QueryRow(ctx, queryGetById, id).Scan(
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)

	if err != nil {
		log.Error(err.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			return post, forumErrors.NewEntityNotExistsError("posts")
		}
	}

	return post, err
}

func (r *PostRepositoryPostgres) GetFromThreadFlat(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetFromThreadFlat",
	})

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id, parent, author, message, is_edited, forum, created").
		From("posts").
		Where("thread = ?", thread)

	if since > 0 {
		if desc {
			queryBuilder = queryBuilder.Where("id < ?", since)
		} else {
			queryBuilder = queryBuilder.Where("id > ?", since)
		}
	}

	if desc {
		queryBuilder = queryBuilder.OrderBy("created DESC, id DESC")
	} else {
		queryBuilder = queryBuilder.OrderBy("created ASC, id ASC")
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

	posts := make([]domain.Post, 0, rows.CommandTag().RowsAffected())
	post := domain.Post{
		Thread: thread,
	}

	for rows.Next() {
		err := rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Created,
		)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepositoryPostgres) GetFromThreadTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetFromThreadTree",
	})
}

func (r *PostRepositoryPostgres) GetFromThreadParentTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetFromThreadParentTree",
	})
}
