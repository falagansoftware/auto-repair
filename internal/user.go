package autorepair

import (
	"context"
	"time"
)

type User struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserCreate struct {
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm-password"`
}

type UserUpdate struct {
	Name     *string
	Surname  *string
	Email    *string
	Active   bool
	Password *string
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	CreateUser(ctx context.Context, user *UserCreate) error
	FindUserByUid(ctx context.Context, uid string) (*User, error)
	FindUsers(ctx context.Context, filters *UserFilters) ([]*User, error)
	FindUsersGlobally(ctx context.Context, search *string) ([]*User, error)
	// UpdateUser(ctx context.Context, user UserUpdate) (*User, error)
	// DeleteUser(ctx context.Context, id string) error
}
