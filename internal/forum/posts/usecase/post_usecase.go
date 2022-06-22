package usecase

import (
	"context"
	"fmt"
	forumsDomain "github.com/rflban/parkmail-dbms/internal/forum/forums/domain"
	"github.com/rflban/parkmail-dbms/internal/forum/posts/domain"
	threadsDomain "github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	usersDomain "github.com/rflban/parkmail-dbms/internal/forum/users/domain"
	"github.com/rflban/parkmail-dbms/internal/pkg/forum/constants"
	"github.com/rflban/parkmail-dbms/pkg/forum/models"
	"github.com/sirupsen/logrus"
	"strconv"
)

type PostRepository interface {
	Create(ctx context.Context, posts []domain.Post) ([]domain.Post, error)
	Patch(ctx context.Context, id int64, message *string) (domain.Post, error)
	GetById(ctx context.Context, id int64) (domain.Post, error)
	GetFromThreadFlat(ctx context.Context, thread string, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadTree(ctx context.Context, thread string, since int64, limit uint64, desc bool) ([]domain.Post, error)
	GetFromThreadParentTree(ctx context.Context, thread string, since int64, limit uint64, desc bool) ([]domain.Post, error)
}

type UserRepository interface {
	GetByNickname(ctx context.Context, nickname string) (usersDomain.User, error)
}

type ThreadRepository interface {
	GetById(ctx context.Context, id int64) (threadsDomain.Thread, error)
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

func (u *PostUseCaseImpl) Create(ctx context.Context, threadSlugOrId string, posts models.Posts) (models.Posts, error) {
	var thread threadsDomain.Thread
	threadId, err := strconv.ParseInt(threadSlugOrId, 10, 64)

	if err != err {
		thread, err = u.threadRepo.GetById(ctx, threadId)
	} else {
		thread, err = u.threadRepo.GetBySlug(ctx, threadSlugOrId)
	}

	threadId32 := int32(threadId)
	toCreate := make([]domain.Post, 0, len(posts))
	for _, post := range posts {
		post.Thread = &threadId32
		post.Forum = &thread.Forum
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
	edited, err := u.postRepo.Patch(ctx, id, message)
	return edited.ToModel(), err
}

func (u *PostUseCaseImpl) GetById(ctx context.Context, id int64) (models.Post, error) {
	obtained, err := u.postRepo.GetById(ctx, id)
	return obtained.ToModel(), err
}

func (u *PostUseCaseImpl) GetDetails(ctx context.Context, id int64, related []string) (models.PostFull, error) {
	log := ctx.Value(constants.UseCaseLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"usecase": "Post",
		"method":  "GetDetails",
	})

	postFull := models.PostFull{}
	post, err := u.postRepo.GetById(ctx, id)
	postModel := post.ToModel()
	postFull.Post = &postModel

	var (
		userObtained   = false
		threadObtained = false
		forumObtained  = false
	)

	if err != nil {
		return postFull, err
	}

	for _, entity := range related {
		switch entity {
		case "user":
			if userObtained {
				break
			}

			user, err := u.userRepo.GetByNickname(ctx, post.Author)
			if err != nil {
				return postFull, err
			}

			userModel := user.ToModel()
			postFull.Author = &userModel

			userObtained = true
		case "thread":
			if threadObtained {
				break
			}

			thread, err := u.threadRepo.GetById(ctx, post.Thread)
			if err != nil {
				return postFull, err
			}

			threadModel := thread.ToModel()
			postFull.Thread = &threadModel

			threadObtained = true
		case "forum":
			if forumObtained {
				break
			}

			forum, err := u.forumRepo.GetBySlug(ctx, post.Forum)
			if err != nil {
				return postFull, err
			}

			forumModel := forum.ToModel()
			postFull.Forum = &forumModel

			forumObtained = true
		default:
			log.Errorf("unexpected related entity: %s", entity)
			return postFull, fmt.Errorf("unexpected related entity: %s", entity)
		}
	}

	return postFull, nil
}

func (u *PostUseCaseImpl) GetFromThread(ctx context.Context, thread string, since int64, limit uint64, desc bool, sort string) (models.Posts, error) {
	log := ctx.Value(constants.UseCaseLogKey).(*logrus.Entry).WithFields(logrus.Fields{
		"usecase": "Post",
		"method":  "GetFromThread",
	})

	var (
		posts []domain.Post
		err   error
	)

	switch sort {
	case "flat":
		posts, err = u.postRepo.GetFromThreadFlat(ctx, thread, since, limit, desc)
	case "tree":
		posts, err = u.postRepo.GetFromThreadTree(ctx, thread, since, limit, desc)
	case "parent_tree":
		posts, err = u.postRepo.GetFromThreadParentTree(ctx, thread, since, limit, desc)
	default:
		log.Errorf("unexpected sort type: %s", sort)
		return nil, fmt.Errorf("unexpected sort type: %s", sort)
	}

	if err != nil {
		return nil, err
	}

	obtained := make(models.Posts, 0, len(posts))
	for _, post := range posts {
		obtained = append(obtained, post.ToModel())
	}

	return obtained, nil
}
