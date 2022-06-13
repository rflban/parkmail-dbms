package repositories

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/service/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	"github.com/sirupsen/logrus"
)

const (
	queryGetStatus   = "SELECT (SELECT COUNT(*) FROM users), (SELECT COUNT(*) FROM forums), (SELECT COUNT(*) FROM threads), (SELECT COUNT(*) FROM posts)"
	queryTruncateAll = "TRUNCATE TABLE users, forums, forums_users, threads, posts, votes CASCADE"
)

type ServiceRepoPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *ServiceRepoPostgres {
	return &ServiceRepoPostgres{
		db: db,
	}
}

func (r *ServiceRepoPostgres) Status(ctx context.Context) (domain.Status, error) {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry)

	status := domain.Status{}
	err := r.db.QueryRow(ctx, queryGetStatus).Scan(
		&status.User,
		&status.Forum,
		&status.Thread,
		&status.Post,
	)

	if err != nil {
		log.Error(err.Error())
	}

	return status, err
}

func (r *ServiceRepoPostgres) Clear(ctx context.Context) error {
	log := ctx.Value(constants.RepoLogKey).(*logrus.Entry)

	_, err := r.db.Exec(ctx, queryTruncateAll)

	if err != nil {
		log.Error(err.Error())
	}

	return err
}
