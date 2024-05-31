package gohtmx

import (
	"context"
	"time"
)

type User struct {
	Uid       string
	Name      string
	Surname   string
	Email     string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserUpdate struct {
	Name    *string
	Surname *string
	Email   *string
	Active  bool
}

type UserFilters struct {
	Uid     *string
	Name    *string
	Surname *string
	Email   *string
	Active  bool
	// Restrict to subset of results.
	Offset int
	Limit  int
	Sort   string
	Order  string
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindUserByUid(ctx context.Context, uid string) (*User, error)
	FindUsers(ctx context.Context, filters *UserFilters) ([]*User, error)
	FindUsersGlobally(ctx context.Context, search *string) ([]*User, error)
	// UpdateUser(ctx context.Context, user UserUpdate) (*User, error)
	// DeleteUser(ctx context.Context, id string) error
}
