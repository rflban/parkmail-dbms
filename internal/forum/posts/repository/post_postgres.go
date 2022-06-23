package repository

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/posts/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	queryGetAfterBatch = `SELECT id, created, batch_idx FROM posts WHERE batch_id = $1 ORDER BY id;`
	queryLastId        = `SELECT MAX(id) FROM posts;`
	queryGetById       = `SELECT parent, author, message, is_edited, forum, thread, created FROM posts WHERE id = $1;`
	queryUpdate        = `UPDATE posts
					SET message = COALESCE(NULLIF(TRIM($2), ''), message), is_edited = ($3 AND message != $2)
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

	batchID := uuid.New()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	now := time.Now()

	copied, err := r.db.CopyFrom(ctx, pgx.Identifier{"posts"}, []string{
		"parent",
		"author",
		"message",
		"forum",
		"thread",
		"created",
		"batch_id",
		"batch_idx",
	}, pgx.CopyFromSlice(len(posts), func(i int) ([]interface{}, error) {
		if posts[i].Created.Equal(time.Time{}) {
			posts[i].Created = now
		}

		post := []interface{}{
			posts[i].Parent,
			posts[i].Author,
			posts[i].Message,
			posts[i].Forum,
			posts[i].Thread,
			posts[i].Created,
			batchID.String(),
			i,
		}

		return post, nil
	}))

	if err != nil {
		log.Error(err.Error())

		if err := tx.Rollback(ctx); err != nil {
			log.Error(err.Error())
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.SQLState() {
			case "23514":
				return nil, forumErrors.NewConflictError(
					pgErr.Message,
				)
			case "23503":
				return nil, forumErrors.NewEntityNotExistsError("users or forum")
			}
		}

		return nil, err
	}

	if int(copied) != len(posts) {
		log.
			WithField("copied", fmt.Sprintf("%d/%d", copied, len(posts))).
			Errorf("Failed bulk insert")

		if err := tx.Rollback(ctx); err != nil {
			log.Error(err.Error())
		}

		return nil, err
	}

	rows, err := r.db.Query(ctx, queryGetAfterBatch, batchID.String())
	if err != nil {
		log.Error(err.Error())
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}

	var batch_idx int
	obtained := make([]domain.Post, 0, rows.CommandTag().RowsAffected())
	post := domain.Post{}

	for rows.Next() {
		err := rows.Scan(
			&post.Id,
			&post.Created,
			&batch_idx,
		)
		if err != nil {
			log.Error(err.Error())
			if err := tx.Rollback(ctx); err != nil {
				return nil, err
			}
			return nil, err
		}
		post.Parent = posts[batch_idx].Parent
		post.Author = posts[batch_idx].Author
		post.Message = posts[batch_idx].Message
		post.IsEdited = posts[batch_idx].IsEdited
		post.Forum = posts[batch_idx].Forum
		post.Thread = posts[batch_idx].Thread

		obtained = append(obtained, post)
	}

	err = tx.Commit(ctx)

	return obtained, err
}

func (r *PostRepositoryPostgres) Patch(ctx context.Context, id int64, message *string) (domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "Patch",
	})

	post := domain.Post{
		Id: id,
	}
	err := r.db.QueryRow(ctx, queryUpdate, id, message, message != nil).Scan(
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

func (r *PostRepositoryPostgres) GetFromThreadFlat(ctx context.Context, thread string, since int64, limit uint64, desc bool) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetFromThreadFlat",
	})

	_, err := strconv.ParseInt(thread, 10, 64)
	threadIsNum := err == nil

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id, parent, author, message, is_edited, forum, thread, created").
		From("posts")

	if threadIsNum {
		queryBuilder = queryBuilder.Where("thread = ?", thread)
	} else {
		queryBuilder = queryBuilder.Where("thread = (SELECT id FROM threads WHERE slug = ?)", thread)
	}

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

	log.Info(query)
	log.Info(args)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	posts := make([]domain.Post, 0, rows.CommandTag().RowsAffected())
	post := domain.Post{}

	for rows.Next() {
		err := rows.Scan(
			&post.Id,
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
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepositoryPostgres) GetFromThreadTree(ctx context.Context, thread string, since int64, limit uint64, desc bool) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetFromThreadTree",
	})

	_, err := strconv.ParseInt(thread, 10, 64)
	threadIsNum := err == nil

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id, parent, author, message, is_edited, forum, thread, created").
		From("posts")

	if threadIsNum {
		queryBuilder = queryBuilder.Where("thread = ?", thread)
	} else {
		queryBuilder = queryBuilder.Where("thread = (SELECT id FROM threads WHERE slug = ?)", thread)
	}

	if since > 0 {
		if desc {
			queryBuilder = queryBuilder.Where("path < (SELECT path FROM posts WHERE id = ?)", since)
		} else {
			queryBuilder = queryBuilder.Where("path > (SELECT path FROM posts WHERE id = ?)", since)
		}
	}

	if desc {
		queryBuilder = queryBuilder.OrderBy("path DESC")
	} else {
		queryBuilder = queryBuilder.OrderBy("path ASC, id ASC")
	}

	if limit > 0 {
		queryBuilder = queryBuilder.Limit(limit)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Info(query)
	log.Info(args)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	posts := make([]domain.Post, 0, rows.CommandTag().RowsAffected())
	post := domain.Post{}

	for rows.Next() {
		err := rows.Scan(
			&post.Id,
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
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepositoryPostgres) GetFromThreadParentTree(ctx context.Context, thread string, since int64, limit uint64, desc bool) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetFromThreadParentTree",
	})

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id, parent, author, message, is_edited, forum, thread, created").
		From("posts")

	_, err := strconv.ParseInt(thread, 10, 64)
	threadIsNum := err == nil

	var threadSqlVal string
	if threadIsNum {
		threadSqlVal = "?"
	} else {
		threadSqlVal = `(SELECT id FROM threads WHERE slug = ?)`
	}

	if since > 0 {
		if desc {
			queryBuilder = queryBuilder.
				Where(
					fmt.Sprintf(`path[1] IN (
							SELECT id
							FROM posts
							WHERE thread = %s AND parent = 0 AND path[1] < (
								SELECT path[1]
								FROM posts
								WHERE id = ?)
							ORDER BY id DESC LIMIT ?)`, threadSqlVal),
					thread,
					since,
					limit,
				).
				OrderBy(`path[1] DESC, path ASC, id ASC`)
		} else {
			queryBuilder = queryBuilder.
				Where(
					fmt.Sprintf(`path[1] IN (
							SELECT id
							FROM posts
							WHERE thread = %s AND parent = 0 AND path[1] > (
								SELECT path[1]
								FROM posts
								WHERE id = ?)
							ORDER BY id ASC LIMIT ?)`, threadSqlVal),
					thread,
					since,
					limit,
				).
				OrderBy(`path ASC, id ASC`)
		}
	} else {
		if desc {
			queryBuilder = queryBuilder.
				Where(
					fmt.Sprintf(`path[1] IN (
							SELECT id
							FROM posts
							WHERE thread = %s AND parent = 0
							ORDER BY id DESC LIMIT ?)`, threadSqlVal),
					thread,
					limit,
				).
				OrderBy(`path[1] DESC, path ASC, id ASC`)
		} else {
			queryBuilder = queryBuilder.
				Where(
					fmt.Sprintf(`path[1] IN (
							SELECT id
							FROM posts
							WHERE thread = %s AND parent = 0
							ORDER BY id ASC LIMIT ?)`, threadSqlVal),
					thread,
					limit,
				).
				OrderBy(`path ASC, id ASC`)
		}
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
	post := domain.Post{}

	for rows.Next() {
		err := rows.Scan(
			&post.Id,
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
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
