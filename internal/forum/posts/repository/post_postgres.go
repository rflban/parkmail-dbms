package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rflban/parkmail-dbms/internal/forum/posts/domain"
)

type PostRepositoryPostgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostRepositoryPostgres {
	return &PostRepositoryPostgres{
		db: db,
	}
}

func (r *PostRepositoryPostgres) CreateAt(ctx context.Context, forum string, thread int64, posts [][]interface{}) (domain.Post, error) {
}

func (r *PostRepositoryPostgres) Patch(ctx context.Context, id int64, message *string) (domain.Post, error) {
}

func (r *PostRepositoryPostgres) GetById(ctx context.Context, id int64) (domain.Post, error) {
}

func (r *PostRepositoryPostgres) GetDetails(ctx context.Context, id int64, related []string) (domain.PostFull, error) {
}

func (r *PostRepositoryPostgres) GetFromThreadFlat(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error) {
}

func (r *PostRepositoryPostgres) GetFromThreadTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error) {
}

func (r *PostRepositoryPostgres) GetFromThreadParentTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error) {
}
