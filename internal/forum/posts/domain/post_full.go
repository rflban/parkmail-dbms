package domain

import (
	forumsDomain "github.com/rflban/parkmail-dbms/internal/forum/forums/domain"
	threadsDomain "github.com/rflban/parkmail-dbms/internal/forum/threads/domain"
	usersDomain "github.com/rflban/parkmail-dbms/internal/forum/users/domain"
)

type PostFull struct {
	Post   *Post
	Author *usersDomain.User
	Thread *threadsDomain.Thread
	Forum  *forumsDomain.Forum
}
