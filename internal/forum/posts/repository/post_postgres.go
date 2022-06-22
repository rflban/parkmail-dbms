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
	"strconv"
	"sync"
)

const (
	queryLastId  = `SELECT MAX(id) FROM posts;`
	queryGetById = `SELECT parent, author, message, is_edited, forum, thread, created FROM posts WHERE id = $1;`
	queryUpdate  = `UPDATE posts
					SET message = COALESCE(NULLIF(TRIM($2), ''), message), is_edited = true
					WHERE id = $1
					RETURNING parent, author, message, is_edited, forum, thread, created;`
)

func getLastId(db *pgxpool.Pool) int64 {
	var id int64;

	err := db.QueryRow(context.Background(), "SELECT MAX(id) FROM posts;").Scan(&id)
	if err != nil {
		id = 0
	}

	return id
}

type PostRepositoryPostgres struct {
	db     *pgxpool.Pool
	lastId int64
	mutex  sync.Mutex
}

func New(db *pgxpool.Pool) *PostRepositoryPostgres {
	return &PostRepositoryPostgres{
		db:     db,
		mutex:  sync.Mutex{},
		lastId: getLastId(db),
	}
}

func (r *PostRepositoryPostgres) Create(ctx context.Context, posts []domain.Post) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "Create",
	})

	r.mutex.Lock()
	lastId := r.lastId
	r.lastId = lastId + int64(len(posts))
	r.mutex.Unlock()

	for i := range posts {
		posts[i].Id = int64(i) + lastId + 1
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	copied, err := r.db.CopyFrom(ctx, pgx.Identifier{"posts"}, []string{
		"id",
		"parent",
		"author",
		"message",
		"forum",
		"thread",
		"created",
	}, pgx.CopyFromSlice(len(posts), func(i int) ([]interface{}, error) {
		return []interface{}{
			posts[i].Id,
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

func (r *PostRepositoryPostgres) GetFromThreadFlat(ctx context.Context, thread string, since int64, limit uint64, desc bool) ([]domain.Post, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "Post",
		"method": "GetFromThreadFlat",
	})

	_, err := strconv.ParseInt(thread, 10, 64)
	threadIsNum := err != nil

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
	threadIsNum := err != nil

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
	threadIsNum := err != nil

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
