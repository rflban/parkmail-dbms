package usecase

import (
	"context"
	forumsDomain "github.com/rflban/parkmail-dbms/internal/forum/forums/domain"
	"github.com/rflban/parkmail-dbms/internal/forum/posts/domain"
	threadsDomain "github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	usersDomain "github.com/rflban/parkmail-dbms/internal/forum/users/domain"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
)

type PostRepository interface {
	Create(ctx context.Context, posts []domain.Post) ([]domain.Post, error)
	Patch(ctx context.Context, id int64, message *string) (domain.Post, error)
	GetById(ctx context.Context, id int64) (domain.Post, error)
	GetFromThreadFlat(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadParentTree(ctx context.Context, thread int64, since int64, limit uint64, desc bool) ([]domain.Post, error)
}

type UserRepository interface {
	GetByNickname(ctx context.Context, nickname string) (usersDomain.User, error)
}

type ThreadRepository interface {
	GetBySlug(ctx context.Context, slug string) (threadsDomain.Thread, error)
}

type ForumRepository interface {
	GetBySlug(ctx context.Context, slug string) (forumsDomain.Forum, error)
}

type PostUseCaseImpl struct {
	postRepo   PostRepository
	userRepo   UserRepository
	threadRepo ThreadRepository
	forumRepo  ForumRepository
}

func New(
	postRepo PostRepository,
	userRepo UserRepository,
	threadRepo ThreadRepository,
	forumRepo ForumRepository,
) *PostUseCaseImpl {
	return &PostUseCaseImpl{
		postRepo:   postRepo,
		userRepo:   userRepo,
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
	}
}

func (u *PostUseCaseImpl) Create(ctx context.Context, posts models.Posts) (models.Posts, error) {
	toCreate := make([]domain.Post, 0, len(posts))
	for _, post := range posts {
		toCreate = append(toCreate, domain.FromModel(post))
	}

	created, err := u.postRepo.Create(ctx, toCreate)

	if err != nil {
		return nil, err
	}

	obtained := make(models.Posts, 0, len(created))
	for _, post := range created {
		obtained = append(obtained, post.ToModel())
	}

	return obtained, nil
}

func (u *PostUseCaseImpl) Patch(ctx context.Context, id int64, message *string) (models.Post, error) {
}

func (u *PostUseCaseImpl) GetById(ctx context.Context, id int64) (models.Post, error) {
}

func (u *PostUseCaseImpl) GetDetails(ctx context.Context, id int64, related string) (models.PostFull, error) {
}

func (u *PostUseCaseImpl) GetFromThread(ctx context.Context, thread int64, since int64, limit uint64, desc bool, sort string) (models.Posts, error) {
	switch sort {

	}
}
