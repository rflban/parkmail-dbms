package main

import (
	"context"
	FasthttpRouter "github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	ForumDelivery "github.com/rflban/parkmail-dbms/internal/forum/forums/delivery"
	ForumRepo "github.com/rflban/parkmail-dbms/internal/forum/forums/repository"
	ForumUseCase "github.com/rflban/parkmail-dbms/internal/forum/forums/usecase"
	PostDelivery "github.com/rflban/parkmail-dbms/internal/forum/posts/delivery"
	PostRepo "github.com/rflban/parkmail-dbms/internal/forum/posts/repository"
	PostUseCase "github.com/rflban/parkmail-dbms/internal/forum/posts/usecase"
	ServiceDelivery "github.com/rflban/parkmail-dbms/internal/forum/service/delivery"
	ServiceRepo "github.com/rflban/parkmail-dbms/internal/forum/service/repository"
	ServiceUseCase "github.com/rflban/parkmail-dbms/internal/forum/service/usecase"
	ThreadDelivery "github.com/rflban/parkmail-dbms/internal/forum/threads/delivery"
	ThreadRepo "github.com/rflban/parkmail-dbms/internal/forum/threads/repository"
	ThreadUseCase "github.com/rflban/parkmail-dbms/internal/forum/threads/usecase"
	UserDelivery "github.com/rflban/parkmail-dbms/internal/forum/users/delivery"
	UserRepo "github.com/rflban/parkmail-dbms/internal/forum/users/repository"
	UserUseCase "github.com/rflban/parkmail-dbms/internal/forum/users/usecase"
	VoteRepo "github.com/rflban/parkmail-dbms/internal/forum/votes/repository"
	VoteUseCase "github.com/rflban/parkmail-dbms/internal/forum/votes/usecase"
)

const prefix = "/api"

func SetupHandlers(ctx context.Context, pool *pgxpool.Pool, router *FasthttpRouter.Router) {
	var (
		serviceRepo = ServiceRepo.New(pool)
		userRepo    = UserRepo.New(pool)
		voteRepo    = VoteRepo.New(pool)
		forumRepo   = ForumRepo.New(pool)
		threadRepo  = ThreadRepo.New(pool)
		postRepo    = PostRepo.New(pool)
	)

	var (
		serviceUseCase = ServiceUseCase.New(serviceRepo)
		userUseCase    = UserUseCase.New(userRepo)
		voteUseCase    = VoteUseCase.New(voteRepo)
		forumUseCase   = ForumUseCase.New(forumRepo)
		threadUseCase  = ThreadUseCase.New(threadRepo)
		postUseCase    = PostUseCase.New(postRepo, userRepo, threadRepo, forumRepo)
	)

	var (
		serviceHandler = ServiceDelivery.New(serviceUseCase)
		userHandler    = UserDelivery.New(userUseCase)
		forumHandler   = ForumDelivery.New(forumUseCase, threadUseCase)
		threadHandler  = ThreadDelivery.New(postUseCase, threadUseCase, voteUseCase)
		postHandler    = PostDelivery.New(postUseCase)
	)

	router.POST(prefix+"/forum/create", forumHandler.Create)
	router.GET(prefix+"/forum/{slug}/details", forumHandler.GetDetails)
	router.POST(prefix+"/forum/{slug}/create", forumHandler.CreateThread)
	router.GET(prefix+"/forum/{slug}/users", forumHandler.GetUsers)
	router.GET(prefix+"/forum/{slug}/threads", forumHandler.GetThreads)

	router.GET(prefix+"/post/{id}/details", postHandler.GetDetails)
	router.POST(prefix+"/post/{id}/details", postHandler.Edit)

	router.POST(prefix+"/service/clear", serviceHandler.Clear)
	router.GET(prefix+"/service/status", serviceHandler.Status)

	router.POST(prefix+"/thread/{slug_or_id}/create", threadHandler.CreatePosts)
	router.GET(prefix+"/thread/{slug_or_id}/details", threadHandler.GetDetails)
	router.POST(prefix+"/thread/{slug_or_id}/details", threadHandler.Edit)
	router.GET(prefix+"/thread/{slug_or_id}/posts", threadHandler.GetPosts)
	router.POST(prefix+"/thread/{slug_or_id}/vote", threadHandler.Vote)

	router.POST(prefix+"/user/{nickname}/create", userHandler.Create)
	router.GET(prefix+"/user/{nickname}/profile", userHandler.GetProfileByNickname)
	router.POST(prefix+"/user/{nickname}/profile", userHandler.EditProfileByNickname)
}
