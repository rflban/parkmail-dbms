package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/user/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	forumErrors "github.com/rflban/parkmail-dbms/internal/pkg/forum/errors"
	"github.com/sirupsen/logrus"
)

const (
	queryCreate               = `INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);`
	queryPatch                = `UPDATE users SET fullname = COALESCE(NULLIF(TRIM($2), ''), fullname), about = COALESCE(NULLIF(TRIM($3), ''), about), email = COALESCE(NULLIF(TRIM($4), ''), email) WHERE nickname = $1 RETURNING fullname, about, email;`
	queryGetByEmail           = `SELECT nickname, fullname, about FROM users WHERE email = $1;`
	queryGetByNickname        = `SELECT fullname, about, email FROM users WHERE nickname = $1;`
	queryGetByEmailOrNickname = `SELECT nickname, fullname, about, email FROM users WHERE email = $1 OR nickname = $2;`
)

type UserRepositoryPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		db: db,
	}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user domain.User) (domain.User, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "User",
		"method": "Create",
	})

	_, err := r.db.Exec(ctx, queryCreate, user.Nickname, user.Fullname, user.About, user.Email)

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

	return user, err
}

func (r *UserRepositoryPostgres) Patch(ctx context.Context, nickname string, partialUser domain.PartialUser) (domain.User, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "User",
		"method": "Patch",
	})

	user := domain.User{Nickname: nickname}

	err := r.db.QueryRow(ctx, queryPatch, nickname, partialUser.Fullname, partialUser.About, partialUser.Email).Scan(
		&user.Fullname,
		&user.About,
		&user.Email,
	)
	if err != nil {
		log.Error(err.Error())

		if err.Error() == pgx.ErrNoRows.Error() {
			return user, forumErrors.NewEntityNotExistsError("users")
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return user, forumErrors.NewUniqueError(
				pgErr.TableName,
				pgErr.ColumnName,
			)
		}
	}

	return user, err
}

func (r *UserRepositoryPostgres) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "User",
		"method": "GetByEmail",
	})

	user := domain.User{Email: email}

	err := r.db.QueryRow(ctx, queryGetByEmail, email).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
	)

	if err != nil {
		log.Error(err.Error())

		if err.Error() == pgx.ErrNoRows.Error() {
			return user, forumErrors.NewEntityNotExistsError("users")
		}
	}

	return user, err
}

func (r *UserRepositoryPostgres) GetByNickname(ctx context.Context, nickname string) (domain.User, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "User",
		"method": "GetByNickname",
	})

	user := domain.User{Nickname: nickname}

	err := r.db.QueryRow(ctx, queryGetByNickname, nickname).Scan(
		&user.Fullname,
		&user.About,
		&user.Email,
	)

	if err != nil {
		log.Error(err.Error())

		if err.Error() == pgx.ErrNoRows.Error() {
			return user, forumErrors.NewEntityNotExistsError("users")
		}
	}

	return user, err
}

func (r *UserRepositoryPostgres) GetByEmailOrNickname(ctx context.Context, email, nickname string) ([]domain.User, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"repo":   "User",
		"method": "GetByEmailOrNickname",
	})

	rows, err := r.db.Query(ctx, queryGetByEmailOrNickname, email, nickname)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	user := domain.User{}

	for rows.Next() {
		err = rows.Scan(
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
