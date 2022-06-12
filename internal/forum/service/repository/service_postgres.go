package repositories

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/service/domain"
)

const (
	queryGetStatus   = "SELECT (SELECT COUNT(*) FROM users), (SELECT COUNT(*) FROM forums), (SELECT COUNT(*) FROM threads), (SELECT COUNT(*) FROM posts)"
	queryTruncateAll = "TRUNCATE TABLE users, forums, forums_users, threads, posts, votes CASCADE"
)

type ServiceRepoPostgres struct {
	db *pgxpool.Pool
}

func NewServiceRepoPostgres(db *pgxpool.Pool) *ServiceRepoPostgres {
	return &ServiceRepoPostgres{
		db: db,
	}
}

func (r *ServiceRepoPostgres) Status(ctx context.Context) (domain.Status, error) {
	status := domain.Status{}
	err := r.db.QueryRow(ctx, queryGetStatus).Scan(
		&status.User,
		&status.Forum,
		&status.Thread,
		&status.Post,
	)

	return status, err
}

func (r *ServiceRepoPostgres) Clear(ctx context.Context) error {
	_, err := r.db.Exec(ctx, queryTruncateAll)
	return err
}
